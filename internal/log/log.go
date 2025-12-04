package log

import (
    stdlog "log"
)

// Init configures the standard library logger with sensible defaults.
// It sets timestamps and short file locations, which aid in tracing
// issues during development and production.
func Init() {
    stdlog.SetFlags(stdlog.LstdFlags | stdlog.Lshortfile)
}