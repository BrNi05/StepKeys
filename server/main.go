package main

import (
	"github.com/getlantern/systray"

	Config "stepkeys/server/config"
	Handler "stepkeys/server/handler"
	Logging "stepkeys/server/logging"
	OS "stepkeys/server/os"
	Tray "stepkeys/server/tray"
	Updater "stepkeys/server/updater"
	Web "stepkeys/server/web"
)

func main() {
	// Compute executable path and directory
	execDir := OS.GetExeDir()

	// Log to file with formatting
	Logging.SetupLogging(execDir)

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
			Logging.WriteToLogFile("Serial listener disabled: " + err.Error())
		}
	}()

	// Start tray menu
	systray.Run(Tray.TrayOnReady, Tray.TrayOnExit)
}
