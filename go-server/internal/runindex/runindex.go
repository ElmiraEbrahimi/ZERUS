package runindex

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	currentIndex int
	currentErr   error
	once         sync.Once
)

// Current returns the run index for the current process, initializing it to 1
// if missing.
func Current(logPath string) (int, error) {
	once.Do(func() {
		currentIndex, currentErr = readOrInitIndex(logPath)
	})
	return currentIndex, currentErr
}

// Advance increments the persisted run index and returns the new value.
func Advance(logPath string) (int, error) {
	indexPath, err := indexFilePath(logPath)
	if err != nil {
		return 0, err
	}
	if err := os.MkdirAll(filepath.Dir(indexPath), 0o755); err != nil {
		return 0, err
	}
	last, ok, err := readIndex(indexPath)
	if err != nil {
		return 0, err
	}
	if !ok {
		last = 0
	}
	next := last + 1
	if err := writeIndexFile(indexPath, next); err != nil {
		return 0, err
	}
	return next, nil
}

// EnsureCSVHeaderWithIndex rewrites the header to include index when missing.
func EnsureCSVHeaderWithIndex(path string, header []string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("csv path is required")
	}
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return nil
	}

	reader := bufio.NewReader(file)
	headerLine, readErr := reader.ReadString('\n')
	if readErr != nil && readErr != io.EOF {
		return readErr
	}
	if strings.TrimSpace(headerLine) == "" {
		return nil
	}

	parsed, err := csv.NewReader(strings.NewReader(strings.TrimRight(headerLine, "\r\n"))).Read()
	if err != nil {
		return err
	}
	if hasIndexColumn(parsed) {
		return nil
	}

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "csv-header-")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}()

	writer := csv.NewWriter(tmp)
	if err := writer.Write(header); err != nil {
		return err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	if _, err := io.Copy(tmp, reader); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	return nil
}

func readOrInitIndex(logPath string) (int, error) {
	indexPath, err := indexFilePath(logPath)
	if err != nil {
		return 0, err
	}
	if err := os.MkdirAll(filepath.Dir(indexPath), 0o755); err != nil {
		return 0, err
	}

	current, ok, err := readIndex(indexPath)
	if err != nil {
		return 0, err
	}
	if ok {
		return current, nil
	}
	current = 1
	if err := writeIndexFile(indexPath, current); err != nil {
		return 0, err
	}
	return current, nil
}

func indexFilePath(logPath string) (string, error) {
	if strings.TrimSpace(logPath) == "" {
		return "", errors.New("log path is required")
	}
	logDir := filepath.Dir(logPath)
	if filepath.Base(logDir) == "memtime" {
		logDir = filepath.Dir(logDir)
	}
	return filepath.Join(logDir, "run_index"), nil
}

func readIndex(path string) (int, bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return 0, false, nil
	}
	parsed, parseErr := strconv.Atoi(trimmed)
	if parseErr != nil {
		return 0, false, parseErr
	}
	return parsed, true, nil
}

func writeIndexFile(path string, index int) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "run_index-")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := fmt.Fprintf(tmp, "%d\n", index); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return nil
}

func hasIndexColumn(cols []string) bool {
	for _, col := range cols {
		if strings.EqualFold(strings.TrimSpace(col), "index") {
			return true
		}
	}
	return false
}
