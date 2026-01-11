package memtime

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"runtime/metrics"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/process"

	"l2alchemy/internal/runindex"
)

var (
	memCSVHeader = []string{
		"index",
		"datetime",
		"source",
		"time",
		"rss_mb",
		"footprint_mb",
		"heap_mb",
		"go_total_mb",
		"rss_delta_mb",
		"footprint_delta_mb",
		"heap_delta_mb",
		"go_total_delta_mb",
	}
	memCSVMu     sync.Mutex
	memCSVWriter *csv.Writer
	memCSVFile   *os.File
	memCSVIndex  int
	memCSVPathMu sync.Mutex
	memCSVPaths  map[string]string
	rssProcOnce  sync.Once
	rssProc      *process.Process
	rssProcErr   error
)

var goMemMetricNames = []string{
	"/memory/classes/heap/objects:bytes",
	"/memory/classes/total:bytes",
}

var goMemSamplePool = sync.Pool{
	New: func() any {
		samples := make([]metrics.Sample, len(goMemMetricNames))
		for i, name := range goMemMetricNames {
			samples[i].Name = name
		}
		return samples
	},
}

// PeakSampleInterval is the default sampling interval for peak memory tracking.
const PeakSampleInterval = time.Millisecond

type osMemSample struct {
	rssBytes       uint64
	rssOK          bool
	footprintBytes uint64
	footprintOK    bool
	pssBytes       uint64
	pssOK          bool
	ussBytes       uint64
	ussOK          bool
}

type goMemSample struct {
	heapObjects uint64
	total       uint64
}

type Sample struct {
	source  string
	start   time.Time
	osBase  osMemSample
	goBase  goMemSample
	hasBase bool
}

func Start(source string) *Sample {
	source = strings.TrimSpace(source)
	if source == "" {
		source = "unknown"
	}
	start := time.Now()
	osSample, _ := osMemorySample()
	goSample := readGoMemSample()
	return &Sample{
		source:  source,
		start:   start,
		osBase:  osSample,
		goBase:  goSample,
		hasBase: true,
	}
}

func (s *Sample) End() {
	if s == nil {
		return
	}
	osSample, _ := osMemorySample()
	goSample := readGoMemSample()
	osBase := osMemSample{}
	goBase := goMemSample{}
	if s.hasBase {
		osBase = s.osBase
		goBase = s.goBase
	}
	logSample(s.start, s.source, time.Since(s.start), osSample, goSample, osBase, goBase)
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
	if err := runindex.EnsureCSVHeaderWithIndex(path, memCSVHeader); err != nil {
		log.Printf("memtime csv: header update failed: %v", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	info, statErr := file.Stat()

	memCSVWriter = csv.NewWriter(file)
	memCSVFile = file
	if idx, err := runindex.Current(path); err != nil {
		log.Printf("memtime csv: index init failed: %v", err)
		memCSVIndex = 0
	} else {
		memCSVIndex = idx
	}
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

func logSample(start time.Time, source string, duration time.Duration, osSample osMemSample, goSample goMemSample, osBase osMemSample, goBase goMemSample) {
	ensureWriter()

	memCSVMu.Lock()
	defer memCSVMu.Unlock()
	if memCSVWriter == nil {
		return
	}

	footprintBytes, footprintOK := combinedFootprintBytesOK(osSample)
	baseFootprintBytes, baseFootprintOK := combinedFootprintBytesOK(osBase)

	row := []string{
		strconv.Itoa(memCSVIndex),
		start.Format("20060102_150405"),
		source,
		fmt.Sprintf("%.6f", duration.Seconds()),
		formatMB(osSample.rssBytes, osSample.rssOK),
		formatMB(footprintBytes, footprintOK),
		formatMB(goSample.heapObjects, true),
		formatMB(goSample.total, true),
		formatDeltaMB(osSample.rssBytes, osSample.rssOK, osBase.rssBytes, osBase.rssOK),
		formatDeltaMB(footprintBytes, footprintOK, baseFootprintBytes, baseFootprintOK),
		formatDeltaMB(goSample.heapObjects, true, goBase.heapObjects, true),
		formatDeltaMB(goSample.total, true, goBase.total, true),
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

func formatMB(bytes uint64, ok bool) string {
	if !ok {
		return ""
	}
	return fmt.Sprintf("%.2f", float64(bytes)/(1024.0*1024.0))
}

func formatDeltaMB(bytes uint64, ok bool, baseBytes uint64, baseOK bool) string {
	if !ok || !baseOK {
		return ""
	}
	if bytes < baseBytes {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", float64(bytes-baseBytes)/(1024.0*1024.0))
}

func combinedFootprintBytesOK(osSample osMemSample) (uint64, bool) {
	if osSample.footprintOK {
		return osSample.footprintBytes, true
	}
	if osSample.pssOK {
		return osSample.pssBytes, true
	}
	if osSample.ussOK {
		return osSample.ussBytes, true
	}
	return 0, false
}

func osMemorySample() (osMemSample, error) {
	sample, err := osFootprintBytes()
	if sample.rssOK {
		return sample, err
	}
	rss, rssErr := rssBytes()
	if rssErr != nil {
		if err == nil {
			err = rssErr
		}
		return sample, err
	}
	sample.rssBytes = rss
	sample.rssOK = true
	return sample, err
}

func rssBytes() (uint64, error) {
	rssProcOnce.Do(func() {
		rssProc, rssProcErr = process.NewProcess(int32(os.Getpid()))
	})
	if rssProcErr != nil {
		return 0, rssProcErr
	}
	info, err := rssProc.MemoryInfo()
	if err != nil {
		return 0, err
	}
	return info.RSS, nil
}

func readGoMemSample() goMemSample {
	samples := goMemSamplePool.Get().([]metrics.Sample)
	metrics.Read(samples)
	out := goMemSample{
		heapObjects: sampleUint64(samples[0].Value),
		total:       sampleUint64(samples[1].Value),
	}
	goMemSamplePool.Put(samples)
	return out
}

func sampleUint64(value metrics.Value) uint64 {
	if value.Kind() != metrics.KindUint64 {
		return 0
	}
	return value.Uint64()
}

func updatePeak(osPeak *osMemSample, goPeak *goMemSample, osSample osMemSample, goSample goMemSample) {
	if osSample.rssOK && (!osPeak.rssOK || osSample.rssBytes > osPeak.rssBytes) {
		osPeak.rssBytes = osSample.rssBytes
		osPeak.rssOK = true
	}
	if osSample.footprintOK && (!osPeak.footprintOK || osSample.footprintBytes > osPeak.footprintBytes) {
		osPeak.footprintBytes = osSample.footprintBytes
		osPeak.footprintOK = true
	}
	if osSample.pssOK && (!osPeak.pssOK || osSample.pssBytes > osPeak.pssBytes) {
		osPeak.pssBytes = osSample.pssBytes
		osPeak.pssOK = true
	}
	if osSample.ussOK && (!osPeak.ussOK || osSample.ussBytes > osPeak.ussBytes) {
		osPeak.ussBytes = osSample.ussBytes
		osPeak.ussOK = true
	}
	if goSample.heapObjects > goPeak.heapObjects {
		goPeak.heapObjects = goSample.heapObjects
	}
	if goSample.total > goPeak.total {
		goPeak.total = goSample.total
	}
}

func peakBytes(osSample osMemSample, goSample goMemSample) uint64 {
	if osSample.rssOK {
		return osSample.rssBytes
	}
	if value, ok := combinedFootprintBytesOK(osSample); ok {
		return value
	}
	if goSample.total != 0 {
		return goSample.total
	}
	return goSample.heapObjects
}

// PeakResult is the result of a peak-memory measurement.
type PeakResult struct {
	PeakBytes uint64
	Time      time.Duration
}

// MeasurePeak samples OS footprint and Go runtime memory classes while fn is running.
// It returns the peak RSS bytes when available, otherwise footprint or Go runtime total bytes.
func MeasurePeak(source string, interval time.Duration, fn func() error) (PeakResult, error) {
	if interval <= 0 {
		interval = PeakSampleInterval
	}
	source = strings.TrimSpace(source)
	if source == "" {
		source = "unknown"
	}

	osSample, _ := osMemorySample()
	goSample := readGoMemSample()

	osBase := osSample
	goBase := goSample

	osPeak := osSample
	goPeak := goSample
	var peakMu sync.Mutex

	start := time.Now()
	done := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)
	// sampler
	go func() {
		defer wg.Done()
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				osSample, _ := osMemorySample()
				goSample := readGoMemSample()
				peakMu.Lock()
				updatePeak(&osPeak, &goPeak, osSample, goSample)
				peakMu.Unlock()
			case <-done:
				return
			}
		}
	}()

	err := fn()
	close(done)
	wg.Wait()

	osSample, _ = osMemorySample()
	goSample = readGoMemSample()
	peakMu.Lock()
	updatePeak(&osPeak, &goPeak, osSample, goSample)
	finalOS := osPeak
	finalGo := goPeak
	finalKey := peakBytes(finalOS, finalGo)
	peakMu.Unlock()

	elapsed := time.Since(start)

	logSample(start, source+" (PEAK)", elapsed, finalOS, finalGo, osBase, goBase)
	if shouldIsolatePeak() {
		debug.FreeOSMemory()
	}

	return PeakResult{PeakBytes: finalKey, Time: elapsed}, err
}

func shouldIsolatePeak() bool {
	val, ok := os.LookupEnv("MEMTIME_ISOLATE")
	if ok {
		switch strings.ToLower(strings.TrimSpace(val)) {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}

	memCSVMu.Lock()
	enabled := memCSVWriter != nil
	memCSVMu.Unlock()
	return enabled
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
