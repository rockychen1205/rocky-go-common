// Package logger provides a simple colored logger
package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

// Logger represents a logger instance
type Logger struct {
	verbose bool
}

// New creates a new Logger instance
func New(verbose bool) *Logger {
	return &Logger{verbose: verbose}
}

// Info prints an info message
func (l *Logger) Info(format string, args ...interface{}) {
	color.Set(color.FgCyan)
	fmt.Printf("[INFO] %s - ", time.Now().Format("15:04:05"))
	color.Unset()
	fmt.Printf(format+"\n", args...)
}

// Success prints a success message
func (l *Logger) Success(format string, args ...interface{}) {
	color.Set(color.FgGreen)
	fmt.Printf("[SUCCESS] %s - ", time.Now().Format("15:04:05"))
	color.Unset()
	fmt.Printf(format+"\n", args...)
}

// Warning prints a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	color.Set(color.FgYellow)
	fmt.Printf("[WARNING] %s - ", time.Now().Format("15:04:05"))
	color.Unset()
	fmt.Printf(format+"\n", args...)
}

// Error prints an error message
func (l *Logger) Error(format string, args ...interface{}) {
	color.Set(color.FgRed)
	fmt.Fprintf(os.Stderr, "[ERROR] %s - ", time.Now().Format("15:04:05"))
	color.Unset()
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// Debug prints a debug message if verbose mode is enabled
func (l *Logger) Debug(format string, args ...interface{}) {
	if !l.verbose {
		return
	}
	color.Set(color.FgMagenta)
	fmt.Printf("[DEBUG] %s - ", time.Now().Format("15:04:05"))
	color.Unset()
	fmt.Printf(format+"\n", args...)
}

// Fatal prints an error message and exits the program
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Error(format, args...)
	os.Exit(1)
}

// NewStdLogger creates a new standard library logger.
func NewStdLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.LstdFlags)
}
