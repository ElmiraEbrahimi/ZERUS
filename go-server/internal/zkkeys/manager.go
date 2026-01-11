package zkkeys

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"syscall"
	"time"
)

// CircuitName is the stable identifier used for key folders and API responses.
type CircuitName string

const (
	CircuitMerkleProof CircuitName = "merkle_proof"
	CircuitVotingBatch CircuitName = "voting_batch"
)

// KeyPaths points to the pk/vk files for a circuit.
type KeyPaths struct {
	PK string `json:"pk"`
	VK string `json:"vk"`
}

// KeyStatus describes whether keys exist on disk.
type KeyStatus struct {
	Circuit     CircuitName `json:"circuit"`
	Paths       KeyPaths    `json:"paths"`
	PKExists    bool        `json:"pk_exists"`
	VKExists    bool        `json:"vk_exists"`
	PKSizeBytes int64       `json:"pk_size_bytes"`
	VKSizeBytes int64       `json:"vk_size_bytes"`
	PKModTime   *time.Time  `json:"pk_mod_time,omitempty"`
	VKModTime   *time.Time  `json:"vk_mod_time,omitempty"`
	PKCreated   *time.Time  `json:"pk_created_time,omitempty"`
	VKCreated   *time.Time  `json:"vk_created_time,omitempty"`
	Ready       bool        `json:"ready"`
}

// KeyGenMetrics summarizes resource usage for compilation + setup.
type KeyGenMetrics struct {
	CompilePeakMB int `json:"compile_peak_mb"`
	CompileTimeMS int `json:"compile_time_ms"`
	SetupPeakMB   int `json:"setup_peak_mb"`
	SetupTimeMS   int `json:"setup_time_ms"`
}

// KeyGenResult is returned by generation calls.
type KeyGenResult struct {
	Circuit   CircuitName   `json:"circuit"`
	Generated bool          `json:"generated"`
	Already   bool          `json:"already_present"`
	Paths     KeyPaths      `json:"paths"`
	Metrics   KeyGenMetrics `json:"metrics"`
}

// Manager owns key locations and provides concurrency-safe generation/status.
type Manager struct {
	baseDir string
	mu      sync.Mutex
	running bool
}

// New constructs a Manager. baseDir should be an absolute or working-directory-relative path.
func New(baseDir string) *Manager {
	return &Manager{baseDir: baseDir}
}

// PathsFor returns the pk/vk filesystem paths for a circuit within baseDir.
func PathsFor(baseDir string, c CircuitName) KeyPaths {
	dir := filepath.Join(baseDir, string(c))
	return KeyPaths{PK: filepath.Join(dir, "pk"), VK: filepath.Join(dir, "vk")}
}

func (m *Manager) pathsFor(c CircuitName) KeyPaths {
	return PathsFor(m.baseDir, c)
}

// Status returns pk/vk existence information for all supported circuits.
func (m *Manager) Status() []KeyStatus {
	circuits := SupportedCircuits()
	out := make([]KeyStatus, 0, len(circuits))
	for _, c := range circuits {
		p := m.pathsFor(c)
		pkMeta := fileMeta(p.PK)
		vkMeta := fileMeta(p.VK)
		out = append(out, KeyStatus{
			Circuit:     c,
			Paths:       p,
			PKExists:    pkMeta.Exists,
			VKExists:    vkMeta.Exists,
			PKSizeBytes: pkMeta.Size,
			VKSizeBytes: vkMeta.Size,
			PKModTime:   pkMeta.ModTime,
			VKModTime:   vkMeta.ModTime,
			PKCreated:   pkMeta.CreatedTime,
			VKCreated:   vkMeta.CreatedTime,
			Ready:       pkMeta.Exists && vkMeta.Exists,
		})
	}
	return out
}

// GenerateAll generates keys for both circuits. If keys already exist for a circuit, it is skipped.
// The operation is serialized; a concurrent caller receives an error.
func (m *Manager) GenerateAll() ([]KeyGenResult, error) {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return nil, errors.New("key generation already in progress")
	}
	m.running = true
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		m.running = false
		m.mu.Unlock()
	}()

	// Ensure base directory exists.
	if err := os.MkdirAll(m.baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("create base dir: %w", err)
	}

	// Generate sequentially to limit peak memory usage.
	circuits := SupportedCircuits()
	results := make([]KeyGenResult, 0, len(circuits))
	for _, c := range circuits {
		r, err := m.Generate(c)
		if err != nil {
			return results, err
		}
		results = append(results, r)
	}
	return results, nil
}

// Generate generates keys for a single circuit. If keys already exist, it returns Already=true.
func (m *Manager) Generate(c CircuitName) (KeyGenResult, error) {
	p := m.pathsFor(c)
	pkMeta, vkMeta := fileMeta(p.PK), fileMeta(p.VK)
	if pkMeta.Exists && vkMeta.Exists {
		return KeyGenResult{Circuit: c, Generated: false, Already: true, Paths: p}, nil
	}

	if err := os.MkdirAll(filepath.Dir(p.PK), 0o755); err != nil {
		return KeyGenResult{}, fmt.Errorf("create circuit dir: %w", err)
	}

	_, _, _,
		compilePeakMB, compileTimeMS,
		setupPeakMB, setupTimeMS,
		err := GenerateKeysToFiles(c, p.PK, p.VK, false)
	if err != nil {
		return KeyGenResult{}, err
	}

	return KeyGenResult{
		Circuit:   c,
		Generated: true,
		Already:   false,
		Paths:     p,
		Metrics: KeyGenMetrics{
			CompilePeakMB: compilePeakMB,
			CompileTimeMS: compileTimeMS,
			SetupPeakMB:   setupPeakMB,
			SetupTimeMS:   setupTimeMS,
		},
	}, nil
}

type meta struct {
	Exists      bool
	Size        int64
	ModTime     *time.Time
	CreatedTime *time.Time
}

func fileMeta(path string) meta {
	st, err := os.Stat(path)
	if err != nil || st.IsDir() {
		return meta{Exists: false}
	}

	mt := st.ModTime()

	// "Created" is best-effort:
	//   - on Linux we expose inode change time (ctime)
	//   - on macOS we expose birth time
	// If neither is available, we fall back to mod time.
	ct := mt
	if sys, ok := st.Sys().(*syscall.Stat_t); ok {
		if t, ok := bestEffortCreatedTime(sys); ok {
			ct = t
		}
	}

	return meta{
		Exists:      true,
		Size:        st.Size(),
		ModTime:     &mt,
		CreatedTime: &ct,
	}
}

// bestEffortCreatedTime attempts to extract a creation-ish timestamp from syscall.Stat_t
// in a portable way (without depending on OS-specific field names at compile time).
//
// Priority order:
//  1. Birth time (macOS / BSD)
//  2. Inode change time (ctime) (Linux)
//
// If neither can be extracted, it returns (zero, false).
func bestEffortCreatedTime(st *syscall.Stat_t) (time.Time, bool) {
	rv := reflect.ValueOf(st)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return time.Time{}, false
	}

	// Common Stat_t field names by OS.
	for _, fname := range []string{"Birthtimespec", "Birthtime", "Ctim", "Ctimespec", "Ctime"} {
		fv := rv.FieldByName(fname)
		if !fv.IsValid() {
			continue
		}
		sec, nsec, ok := timespecToUnix(fv)
		if !ok {
			continue
		}
		if sec != 0 {
			return time.Unix(sec, nsec), true
		}
	}
	return time.Time{}, false
}

func timespecToUnix(v reflect.Value) (sec int64, nsec int64, ok bool) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 0, 0, false
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return 0, 0, false
	}

	// Typical names: Sec/Nsec (Go timespec), Tv_sec/Tv_nsec (C), or similar.
	secField := v.FieldByName("Sec")
	nsecField := v.FieldByName("Nsec")
	if secField.IsValid() && nsecField.IsValid() {
		return secField.Int(), nsecField.Int(), true
	}
	secField = v.FieldByName("Tv_sec")
	nsecField = v.FieldByName("Tv_nsec")
	if secField.IsValid() && nsecField.IsValid() {
		return secField.Int(), nsecField.Int(), true
	}
	return 0, 0, false
}
