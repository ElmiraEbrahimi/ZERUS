package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Setup configures the standard logger to write to stdout and a log file.
// The caller is responsible for closing the returned file.
func Setup(logPath string) (*os.File, error) {
	if logPath == "" {
		logPath = "server.log"
	}
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	return logFile, nil
}
