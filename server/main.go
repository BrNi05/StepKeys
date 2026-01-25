package main

import (
	"log"

	"github.com/getlantern/systray"

	Config "stepkeys/server/config"
	Handler "stepkeys/server/handler"
	Logging "stepkeys/server/logging"
	Tray "stepkeys/server/tray"
	Updater "stepkeys/server/updater"
	Web "stepkeys/server/web"
)

func main() {
	// Log to file with formatting
	Logging.SetupLogging()

	// Load .env or use defaults
	baudRate, serialPort := Config.LoadEnv()

	// Load app and pedal map config
	Config.LoadConfig()

	// Check for updates
	Updater.CheckForUpdates()

	// Register API docs route
	Web.ServeApiDocs()

	// Register API routes
	Web.RegisterAPI()

	// Start webGUI
	go Web.StartGUI()

	// Start serial listener
	// If StepKeys is not enabled, the listener will not process any events
	go func() {
		if err := Handler.ListenSerial(baudRate, serialPort); err != nil {
			log.Println("Serial listener disabled:", err)
		}
	}()

	// Start tray menu
	systray.Run(Tray.TrayOnReady, Tray.TrayOnExit)
}
