package logging

import (
	"log"
	"os"
)

func SetupLogging() {
	logFile, err := os.OpenFile("stepkeys.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logFile)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
