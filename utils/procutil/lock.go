package procutil

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Lock manages a file-based process lock to prevent multiple instances.
type Lock struct {
	filePath string
	logger   *log.Logger
}

// New creates a new Lock for a given application name.
// It uses a default logger if none is provided.
func New(appName string, logger *log.Logger) *Lock {
	if logger == nil {
		logger = log.New(os.Stderr, "[ProcLock] ", log.LstdFlags)
	}
	return &Lock{
		filePath: getLockFilePath(appName),
		logger:   logger,
	}
}

// Acquire attempts to acquire the process lock.
func (l *Lock) Acquire() error {
	// Check if lock file exists and if the process is still running.
	if _, err := os.Stat(l.filePath); err == nil {
		data, err := os.ReadFile(l.filePath)
		if err == nil {
			pidStr := strings.TrimSpace(string(data))
			if pid, err := strconv.Atoi(pidStr); err == nil {
				if l.isProcessRunning(pid) {
					return fmt.Errorf("another instance (PID: %d) is already running, lock file: %s", pid, l.filePath)
				}
				l.logger.Printf("Removing stale lock file for PID %d.", pid)
				os.Remove(l.filePath)
			}
		}
	}

	// Create a new lock file.
	file, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("another instance is running (lock file exists: %s)", l.filePath)
		}
		return fmt.Errorf("could not create lock file: %w", err)
	}
	defer file.Close()

	// Write current PID to the lock file.
	currentPID := os.Getpid()
	if _, err := file.WriteString(fmt.Sprintf("%d\n", currentPID)); err != nil {
		os.Remove(l.filePath) // Clean up on failure.
		return fmt.Errorf("failed to write PID to lock file: %w", err)
	}

	if err := file.Sync(); err != nil {
		os.Remove(l.filePath) // Clean up on failure.
		return fmt.Errorf("failed to sync lock file: %w", err)
	}

	l.logger.Printf("Process lock acquired (PID: %d, File: %s)", currentPID, l.filePath)
	return nil
}

// Release removes the process lock file.
func (l *Lock) Release() {
	if err := os.Remove(l.filePath); err != nil {
		l.logger.Printf("Warning: could not remove lock file %s: %v", l.filePath, err)
	} else {
		l.logger.Printf("Process lock released (File: %s)", l.filePath)
	}
}

// getLockFilePath returns a platform-specific path for the lock file.
func getLockFilePath(appName string) string {
	if runtime.GOOS == "windows" {
		tempDir := os.Getenv("TEMP")
		if tempDir == "" {
			tempDir = os.Getenv("TMP")
		}
		if tempDir == "" {
			tempDir = "C:\\Windows\\Temp"
		}
		return fmt.Sprintf("%s\\%s.lock", tempDir, appName)
	}
	// For Unix-like systems.
	return fmt.Sprintf("/tmp/%s.lock", appName)
}

// isProcessRunning checks if a process with the given PID is currently running.
func (l *Lock) isProcessRunning(pid int) bool {
	if runtime.GOOS == "windows" {
		return l.isProcessRunningWindows(pid)
	}
	return l.isProcessRunningUnix(pid)
}

// isProcessRunningWindows checks for a process on Windows.
// Note: This is the original logic from the app. A more robust check might require platform-specific APIs (e.g., syscall).
func (l *Lock) isProcessRunningWindows(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Windows, os.FindProcess almost always succeeds. This check is not highly reliable
	// but is retained to match original behavior without adding new dependencies.
	return process != nil
}

// isProcessRunningUnix checks for a process on Unix-like systems (specifically those with /proc).
func (l *Lock) isProcessRunningUnix(pid int) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	return err == nil
}
