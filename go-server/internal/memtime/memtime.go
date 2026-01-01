package memtime

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
)

var (
	memCSVHeader = []string{"datetime", "source", "time", "memory"}
	memCSVMu     sync.Mutex
	memCSVWriter *csv.Writer
	memCSVFile   *os.File
	memCSVPathMu sync.Mutex
	memCSVPaths  map[string]string
)

type Sample struct {
	source string
	start  time.Time
	before runtime.MemStats
}

func Start(source string) *Sample {
	source = strings.TrimSpace(source)
	if source == "" {
		source = "unknown"
	}
	var before runtime.MemStats
	runtime.ReadMemStats(&before)
	return &Sample{
		source: source,
		start:  time.Now(),
		before: before,
	}
}

func (s *Sample) End() {
	if s == nil {
		return
	}
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	usedBytes := uint64(0)
	if after.TotalAlloc > s.before.TotalAlloc {
		usedBytes = after.TotalAlloc - s.before.TotalAlloc
	}
	logSample(s.start, s.source, time.Since(s.start), usedBytes)
}

// SetupMemTimeCSV configures CSV logging for memory/time samples.
// The caller is responsible for calling CloseMemTimeCSV.
func SetupMemTimeCSV(path string) (*os.File, error) {
	if strings.TrimSpace(path) == "" {
		return nil, nil
	}

	memCSVMu.Lock()
	defer memCSVMu.Unlock()
	if memCSVWriter != nil {
		return memCSVFile, nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	info, statErr := file.Stat()

	memCSVWriter = csv.NewWriter(file)
	memCSVFile = file
	if statErr == nil && info.Size() == 0 {
		if err := memCSVWriter.Write(memCSVHeader); err != nil {
			log.Printf("memtime csv: header write failed: %v", err)
		}
		memCSVWriter.Flush()
	}

	return file, nil
}

// CloseMemTimeCSV flushes and closes the CSV logger if configured.
func CloseMemTimeCSV() {
	memCSVMu.Lock()
	if memCSVWriter != nil {
		memCSVWriter.Flush()
	}
	if memCSVFile != nil {
		if err := memCSVFile.Close(); err != nil {
			log.Printf("memtime csv: close failed: %v", err)
		}
	}
	memCSVWriter = nil
	memCSVFile = nil
	memCSVMu.Unlock()
}

// SetupFromEnv attempts to load .env and configure the CSV logger from NODE_COUNT and BATCH_SIZE.
func SetupFromEnv() (*os.File, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	nodeCount, err := envInt("NODE_COUNT")
	if err != nil {
		return nil, err
	}
	batchSize, err := envInt("BATCH_SIZE")
	if err != nil {
		return nil, err
	}

	return SetupMemTimeCSV(ResolveMemTimeCSVPath(nodeCount, batchSize))
}

// ResolveMemTimeCSVPath builds the log path for a node/batch pair.
func ResolveMemTimeCSVPath(nodeCount, batchSize int) string {
	cacheKey := fmt.Sprintf("n%d_b%d", nodeCount, batchSize)
	memCSVPathMu.Lock()
	if memCSVPaths != nil {
		if cached, ok := memCSVPaths[cacheKey]; ok {
			memCSVPathMu.Unlock()
			return cached
		}
	}
	memCSVPathMu.Unlock()

	dir := resolveMemTimeDir()
	filename := fmt.Sprintf("n%d_b%d.csv", nodeCount, batchSize)
	path := filepath.Join(dir, filename)

	memCSVPathMu.Lock()
	if memCSVPaths == nil {
		memCSVPaths = make(map[string]string)
	}
	memCSVPaths[cacheKey] = path
	memCSVPathMu.Unlock()

	return path
}

func logSample(start time.Time, source string, duration time.Duration, usedBytes uint64) {
	ensureWriter()

	memCSVMu.Lock()
	defer memCSVMu.Unlock()
	if memCSVWriter == nil {
		return
	}

	row := []string{
		start.Format("20060102_150405"),
		source,
		fmt.Sprintf("%.6f", duration.Seconds()),
		fmt.Sprintf("%.2f", float64(usedBytes)/(1024.0*1024.0)),
	}

	if err := memCSVWriter.Write(row); err != nil {
		log.Printf("memtime csv: write failed: %v", err)
		return
	}
	memCSVWriter.Flush()
	if err := memCSVWriter.Error(); err != nil {
		log.Printf("memtime csv: flush failed: %v", err)
	}
}

// PeakResult is the result of a peak-memory measurement.
type PeakResult struct {
	PeakBytes uint64
	Time      time.Duration
}

// MeasurePeak samples runtime.MemStats.Alloc while fn is running.
// It returns (peakAllocDuringFn - baselineAllocBeforeFn).
func MeasurePeak(source string, interval time.Duration, fn func() error) (PeakResult, error) {
	if interval <= 0 {
		interval = 5 * time.Millisecond
	}
	source = strings.TrimSpace(source)
	if source == "" {
		source = "unknown"
	}

	// Reduce noise
	runtime.GC()

	var base runtime.MemStats
	runtime.ReadMemStats(&base)

	var peak uint64
	peak = base.Alloc

	start := time.Now()
	done := make(chan struct{})

	// sampler
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				var ms runtime.MemStats
				runtime.ReadMemStats(&ms)
				// update peak
				for {
					old := atomic.LoadUint64(&peak)
					if ms.Alloc <= old {
						break
					}
					if atomic.CompareAndSwapUint64(&peak, old, ms.Alloc) {
						break
					}
				}
			case <-done:
				return
			}
		}
	}()

	err := fn()
	close(done)

	elapsed := time.Since(start)
	p := atomic.LoadUint64(&peak)

	used := uint64(0)
	if p > base.Alloc {
		used = p - base.Alloc
	}

	// Log peak as MB to your CSV (same CSV system you already have)
	logSample(start, source+" (PEAK)", elapsed, used)

	return PeakResult{PeakBytes: used, Time: elapsed}, err
}

func BytesToMB(b uint64) int {
	return int(b / (1024 * 1024))
}

func ensureWriter() {
	memCSVMu.Lock()
	ready := memCSVWriter != nil
	memCSVMu.Unlock()
	if ready {
		return
	}
	if !autoSetupEnabled() {
		return
	}
	if _, err := SetupFromEnv(); err != nil {
		log.Printf("memtime csv: setup failed: %v", err)
	}
}

func autoSetupEnabled() bool {
	val, ok := os.LookupEnv("MEMTIME_AUTO_SETUP")
	if !ok {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func envInt(key string) (int, error) {
	val, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(val) == "" {
		return 0, fmt.Errorf("%s is required for memtime logging", key)
	}
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q: %w", key, val, err)
	}
	return n, nil
}

func resolveMemTimeDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return filepath.Join("logs", "memtime")
	}
	if st, err := os.Stat(filepath.Join(wd, "go-server")); err == nil && st.IsDir() {
		return filepath.Join(wd, "go-server", "logs", "memtime")
	}
	return filepath.Join(wd, "logs", "memtime")
}
