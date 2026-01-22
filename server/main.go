package main

import (
	"log"

	"github.com/getlantern/systray"

	Config "stepkeys/server/config"
	Logging "stepkeys/server/logging"
	OS "stepkeys/server/os"
	Serial "stepkeys/server/serial"
	Web "stepkeys/server/web"
)

func main() {
	// Log to file with formatting
	Logging.SetupLogging()

	// Load .env or use defaults
	Config.LoadEnv()

	// Load app and pedal map config
	Config.LoadConfig()

	// Register API docs route
	Web.ServeApiDocs()

	// Register API routes
	Web.RegisterAPI()

	// Start webGUI
	go Web.StartGUI()

	// Start serial listener
	// If StepKeys is not enabled, the listener will not process any events
	go func() {
		if err := Serial.ListenSerial(); err != nil {
			log.Println("Serial listener disabled:", err)
		}
	}()

	// Start tray menu
	systray.Run(OS.TrayOnReady, OS.TrayOnExit)
}
