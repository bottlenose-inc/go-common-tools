package logger

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	Name       string
	Hostname   string
	Pid        int
	LogLevel   int
	file       *os.File
	writer     io.Writer
	isBuffered bool
	lock       sync.Mutex
}

const (
	TraceLevel          int = 10
	DebugLevel          int = 20
	InfoLevel           int = 30
	WarnLevel           int = 40
	ErrorLevel          int = 50
	FatalLevel          int = 60
	BunyanSyntaxVersion int = 0
)

// Returns a fully configured Logger
func NewLogger(name string, args ...string) (*Logger, error) {
	file, err := parseArgs(args...)
	if err != nil {
		return nil, err
	}
	return newLogger(name, file, file), nil
}

// Returns a fully configured Buffered Logger
func NewBufferedLogger(name string, bufSize int, args ...string) (*Logger, error) {
	file, err := parseArgs(args...)
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriterSize(file, bufSize)
	logger := newLogger(name, writer, file)
	logger.isBuffered = true
	return logger, nil
}

// Set LogLevel, only supports the levels defined as consts above
// Defaults to TraceLevel (all logs will be written)
func (logger *Logger) SetLogLevel(level string) {
	switch level {
	case "fatal":
		logger.LogLevel = FatalLevel
	case "error":
		logger.LogLevel = ErrorLevel
	case "warn":
		logger.LogLevel = WarnLevel
	case "info":
		logger.LogLevel = InfoLevel
	case "debug":
		logger.LogLevel = DebugLevel
	default:
		logger.LogLevel = TraceLevel
	}
}

func newLogger(name string, writer io.Writer, file *os.File) *Logger {
	logger := new(Logger)
	logger.Name = strings.TrimSpace(name)
	logger.Hostname, _ = os.Hostname()
	logger.Pid = os.Getpid()
	logger.file = file
	logger.writer = writer
	return logger
}

// Returns os.File based on args
func parseArgs(args ...string) (*os.File, error) {
	if args != nil { // We only care about args[0], but using ...string allows args to be omitted
		path := strings.Replace(strings.TrimSpace(args[0]), "\\", "/", -1)
		// Creates path to log file if it does not already exist
		if strings.Contains(path, "/") {
			if err := os.MkdirAll(path[0:strings.LastIndex(path, "/")], 0777); err != nil {
				return nil, err
			}
		}
		// Open log file
		return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	} else {
		return os.Stdout, nil
	}
}

// Required for expected output if using a Buffered Logger, recommended otherwise
func (logger *Logger) Close() (flushErr error, closeErr error) {
	// Protect access to writer & file
	logger.lock.Lock()
	defer logger.lock.Unlock()

	// Flush buffer (if buffered logger) and close file
	if logger.isBuffered {
		flushErr = logger.writer.(*bufio.Writer).Flush()
	}
	if logger.file != os.Stdout {
		closeErr = logger.file.Close()
	}
	return flushErr, closeErr
}

// Log outputs a JSON-ified log to the configured destination
func (logger *Logger) Log(msg string, level int, extras ...map[string]string) error {
	// Create initial log entry map
	logEntry := map[string]interface{}{
		"hostname": logger.Hostname,
		"level":    level,
		"msg":      msg,
		"name":     logger.Name,
		"pid":      logger.Pid,
		"time":     strings.Replace(time.Now().String()[:23], " ", "T", 1) + "Z", // time in bunyan's format
		"v":        BunyanSyntaxVersion,
	}

	// Add extras to log entry if provided
	if extras != nil {
		for _, extra := range extras {
			for field, value := range extra {
				logEntry[field] = value
			}
		}
	}

	// Protect access to writer
	logger.lock.Lock()
	defer logger.lock.Unlock()

	// Marshal log entry to JSON, or log error
	if logJson, err := json.Marshal(logEntry); err != nil {
		io.WriteString(logger.writer, fmt.Sprintf("Error marshalling log entry JSON: %s", err.Error()))
		return err
	} else {
		// Write log entry
		_, err := io.WriteString(logger.writer, string(logJson)+"\n")
		if err != nil {
			logger.writer = os.Stdout
			logger.Error(fmt.Sprintf("Error writing to log: %s", err.Error()))
			return err
		}
	}
	return nil
}

// Trace writes a log at TraceLevel
func (logger *Logger) Trace(msg string, extras ...map[string]string) error {
	if TraceLevel >= logger.LogLevel {
		return logger.Log(msg, TraceLevel, extras...)
	}
	return nil
}

// Debug writes a log at DebugLevel
func (logger *Logger) Debug(msg string, extras ...map[string]string) error {
	if DebugLevel >= logger.LogLevel {
		return logger.Log(msg, DebugLevel, extras...)
	}
	return nil
}

// Info writes a log at InfoLevel
func (logger *Logger) Info(msg string, extras ...map[string]string) error {
	if InfoLevel >= logger.LogLevel {
		return logger.Log(msg, InfoLevel, extras...)
	}
	return nil
}

// Warning writes a log at WarnLevel
func (logger *Logger) Warning(msg string, extras ...map[string]string) error {
	if WarnLevel >= logger.LogLevel {
		return logger.Log(msg, WarnLevel, extras...)
	}
	return nil
}

// Error writes a log at ErrorLevel
func (logger *Logger) Error(msg string, extras ...map[string]string) error {
	if ErrorLevel >= logger.LogLevel {
		return logger.Log(msg, ErrorLevel, extras...)
	}
	return nil
}

// Fatal writes a log at FatalLevel
func (logger *Logger) Fatal(msg string, extras ...map[string]string) error {
	if FatalLevel >= logger.LogLevel {
		return logger.Log(msg, FatalLevel, extras...)
	}
	return nil
}
