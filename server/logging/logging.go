package logging

import (
	"log"
	"os"
	"path/filepath"

	OS "stepkeys/server/os"
)

func SetupLogging() {
	logFilePath := filepath.Join(OS.GetExecDir(), "stepkeys.log")

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logFile)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Start logging session
	log.Println("\nStepKeys started.")
}
