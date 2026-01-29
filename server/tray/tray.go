package tray

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/pkg/browser"

	Config "stepkeys/server/config"
	Log "stepkeys/server/logging"
	Updater "stepkeys/server/updater"
)

//go:embed assets/icon.ico
var iconBytes []byte

func TrayOnReady() {
	Log.WriteToLogFile("StepKeys tray menu is initialized.")

	systray.SetTitle("")
	systray.SetTooltip("StepKeys Server")
	systray.SetIcon(iconBytes)

	menuOpen := systray.AddMenuItem("Open", "")
	systray.AddSeparator()
	menuEnabled := systray.AddMenuItemCheckbox("Enabled", "", Config.IsEnabled())
	menuStart := systray.AddMenuItemCheckbox("Start on boot", "", Config.IsStartOnBootEnabled())
	systray.AddSeparator()
	menuApiDocs := systray.AddMenuItem("API Docs", "")
	menuDocs := systray.AddMenuItem("Docs", "")
	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("Quit", "")

	// Indicate if update is available
	if Updater.UpdateAvailable() {
		menuOpen.SetTitle("Open (Update Available)")
	}

	// Initial state
	if Config.IsEnabled() {
		menuEnabled.Check()
	} else {
		menuEnabled.Uncheck()
	}

	if Config.IsStartOnBootEnabled() {
		menuStart.Check()
	} else {
		menuStart.Uncheck()
	}

	// Tray event loop
	go func() {
		for {
			select {
			// Open web GUI
			case <-menuOpen.ClickedCh:
				menuOpen.Uncheck()
				openBrowser(fmt.Sprintf("http://localhost:%d", Config.GetWebPort()))

			// Toggle if enabled
			case <-menuEnabled.ClickedCh:
				Config.ToggleEnabled()
				if Config.IsEnabled() {
					menuEnabled.Check()
				} else {
					menuEnabled.Uncheck()
				}

			// Toggle start on boot
			case <-menuStart.ClickedCh:
				Config.ToggleStartOnBoot()
				if Config.IsStartOnBootEnabled() {
					menuStart.Check()
				} else {
					menuStart.Uncheck()
				}

			// Open API Docs page
			case <-menuApiDocs.ClickedCh:
				menuApiDocs.Uncheck()
				openBrowser(fmt.Sprintf("http://localhost:%d/api/docs", Config.GetWebPort()))

			// Open StepKeys GitHub page
			case <-menuDocs.ClickedCh:
				menuDocs.Uncheck()
				openBrowser("https://github.com/BrNi05/StepKeys?tab=readme-ov-file")

			// Quit StepKeys
			case <-menuQuit.ClickedCh:
				menuQuit.Uncheck()
				systray.Quit()
				os.Exit(0)
			}
		}
	}()

	// API event loop
	go func() {
		for {
			select {
			case val := <-Config.EnabledChanged():
				if val {
					menuEnabled.Check()
				} else {
					menuEnabled.Uncheck()
				}

			case val := <-Config.StartOnBootChanged():
				if val {
					menuStart.Check()
				} else {
					menuStart.Uncheck()
				}
			}
		}
	}()
}

func TrayOnExit() { Log.WriteToLogFile("StepKeys server is shutting down.") }

// Browser open helper
func openBrowser(url string) {
	err := browser.OpenURL(url)
	if err != nil {
		Log.WriteToLogFile(fmt.Sprintf("Failed to open browser: %v.", err))
	} else {
		Log.WriteToLogFile(fmt.Sprintf("Opened browser to %s.", url))
	}
}
