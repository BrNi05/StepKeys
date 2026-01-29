package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const sessionStartMarker = "StepKeys started."

// List of all log entries that was written to file
var logs []string

func SetupLogging(execDir string) {
	logFilePath := filepath.Join(execDir, "stepkeys.log")

	// Check if log file is too big and reset if needed
	resetted := resetLogFile(logFilePath)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logFile)

	log.SetFlags(log.LstdFlags)

	// Start logging session
	WriteToLogFile("\n" + sessionStartMarker)

	if resetted {
		WriteToLogFile("Old log file was reset and a new was opened.")
	}
}

func WriteToLogFile(message string) {
	log.Println(message)

	// Format
	message = strings.TrimLeft(message, "\n")
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	line := fmt.Sprintf("%s %s", timestamp, message)

	broadcastLog(line)        // broadcast to WebSocket clients (GUI)
	logs = append(logs, line) // store in memory, used by API
}

func ErrorToLogFile(message string) {
	WriteToLogFile(message)
	os.Exit(1)
}

// Called once a log file gets too big (4 MB)
// Deletes the log file and starts a new one
func resetLogFile(logFilePath string) bool {
	if !isLogFileTooBig(logFilePath, 4) {
		return false
	}

	err := os.Remove(logFilePath)
	if err != nil {
		return false // silent error
	}

	return true
}

// Checks if the log file exceeds the given size in megabytes
func isLogFileTooBig(logFilePath string, maxSizeMegabytes int64) bool {
	fileInfo, err := os.Stat(logFilePath)
	if err != nil {
		return false // silent error
	}

	fileSize := fileInfo.Size()
	maxSizeBytes := maxSizeMegabytes * 1024 * 1024

	return fileSize > maxSizeBytes
}

// Reads the whole log file and returns only the lines from the current session
// Used by the API to serve log content (on initial GUI load)
func ReadCurrentSessionLogs() []string {
	return logs
}
