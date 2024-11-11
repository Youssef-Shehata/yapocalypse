package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}
type logLevel string

// Define constants for valid log levels
const (
	ERROR logLevel = "ERROR"
	INFO  logLevel = "INFO"
)

func validateLogLevel(level logLevel) bool {
	switch level {
	case ERROR, INFO:
		return true
	default:
		return false
	}
}

func NewLogger(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &Logger{
		file: file,
	}, nil
}

func (l *Logger) Log(lvl logLevel , err error ) {
	if !validateLogLevel(lvl) {
		if os.Getenv("PLATFORM") == "dev" {
			log.Fatalf("Invalid log level: %s", lvl)
		}
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(" [%s] %v : %v\n", timestamp, lvl, err)

	fmt.Println(msg)

	if _, writeErr := l.file.WriteString(msg); writeErr != nil {
		log.Printf("Failed to write to log file: %v\n", writeErr)
	}
}

func (l *Logger) Close() error {
	return l.file.Close()
}

