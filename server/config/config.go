package config

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	Handler "stepkeys/server/handler"
	Log "stepkeys/server/logging"
	OS "stepkeys/server/os"
	. "stepkeys/server/pedal"
)

// Key-value map of pedal IDs to actions
var (
	pedalMap   = make(PedalMap)
	pedalMapMu sync.RWMutex
)

// State change events
var (
	enabledChanged     = make(chan bool, 4)
	startOnBootChanged = make(chan bool, 4)
)

// Channel to notify enabled state changes
func EnabledChanged() <-chan bool {
	return enabledChanged
}

// Channel to notify start on boot state changes
func StartOnBootChanged() <-chan bool {
	return startOnBootChanged
}

// App config structure
type AppConfig struct {
	WebPort     int  `json:"webPort"`
	StartOnBoot bool `json:"startOnBoot"`
	Enabled     bool `json:"enabled"`
}

var (
	appConfig   AppConfig
	appConfigMu sync.RWMutex

	defaultPort = 18000 // preferred web server port
)

// Path to stepkeys executable
// Used to resolve relative paths
var execPath string

// File paths
var appConfigFilePath string
var pedalConfigFilePath string

func initConfigFiles() {
	// The path to the executable (can be a symlink)
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Resolve symlinks
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)
	execPath = exePath

	appConfigFilePath = filepath.Join(exeDir, "config.json")
	pedalConfigFilePath = filepath.Join(exeDir, "pedals.json")
}

// Load config from file
func LoadConfig() {
	initConfigFiles()
	Log.WriteToLogFile("Loading app and pedal config.")

	// Load app config
	// Start on boot state is not enforced here, only on runtime toggle events
	if data, err := os.ReadFile(appConfigFilePath); err == nil {
		if err := json.Unmarshal(data, &appConfig); err != nil {
			Log.WriteToLogFile("Error parsing app config, using defaults: " + err.Error())
			appConfig = AppConfig{WebPort: defaultPort, StartOnBoot: false, Enabled: false}
			saveAppConfig()
		}
	} else {
		Log.WriteToLogFile("App config not found, creating default")
		appConfig = AppConfig{WebPort: defaultPort, StartOnBoot: false, Enabled: false}
		saveAppConfig()
	}

	// Load pedal config
	if data, err := os.ReadFile(pedalConfigFilePath); err == nil {
		if err := json.Unmarshal(data, &pedalMap); err != nil {
			Log.WriteToLogFile("Error parsing pedal config, disabling pedals: " + err.Error())
			pedalMap = make(PedalMap)
			if IsEnabled() {
				ToggleEnabled() // disable if StepKeys was enabled
			}
		}
	} else {
		Log.WriteToLogFile("Pedal config not found, disabling pedals")
		pedalMap = make(PedalMap)
		if IsEnabled() {
			ToggleEnabled() // disable if StepKeys was enabled
		}
	}

	// Sync handler copies
	Handler.UpdatePedalMap(GetPedalMap())
	Handler.UpdateEnabled(IsEnabled())
}

// Save config data to file
// Used for initial config file generation and tray menu triggered overrides
func saveAppConfig() {
	data, _ := json.MarshalIndent(appConfig, "", "  ")
	if err := os.WriteFile(appConfigFilePath, data, 0644); err != nil {
		Log.WriteToLogFile("Failed to save app config: " + err.Error())
	} else {
		Log.WriteToLogFile("App config saved.")
	}
}

// Returns the enabled state
func IsEnabled() bool {
	appConfigMu.RLock()
	defer appConfigMu.RUnlock()
	return appConfig.Enabled
}

// Flips the enabled state
func ToggleEnabled() {
	appConfigMu.Lock()

	// Pre-enable checks
	if !appConfig.Enabled && len(pedalMap) == 0 {
		appConfigMu.Unlock()
		return
	}

	appConfig.Enabled = !appConfig.Enabled
	stateSnapshot := appConfig.Enabled
	Handler.UpdateEnabled(stateSnapshot) // update handler copy

	saveAppConfig() // make changes persistent

	Log.WriteToLogFile("StepKeys enable state was changed: " + strconv.FormatBool(appConfig.Enabled))

	appConfigMu.Unlock()

	// Notify possible event subscribers
	select {
	case enabledChanged <- stateSnapshot:
	default:
	}
}

// Returns if start on boot is enabled
// Used by the tray menu
func IsStartOnBootEnabled() bool {
	appConfigMu.RLock()
	defer appConfigMu.RUnlock()
	return appConfig.StartOnBoot
}

// Flips the start on boot state
// Used by the tray menu
func ToggleStartOnBoot() {
	appConfigMu.Lock()
	appConfig.StartOnBoot = !appConfig.StartOnBoot
	startOnBootSnapshot := appConfig.StartOnBoot

	// Process changes
	if appConfig.StartOnBoot {
		if err := OS.EnableStartOnBoot(); err != nil {
			Log.WriteToLogFile("Failed to enable OS-level autostart: " + err.Error())
			return
		}
		Log.WriteToLogFile("Autostart enabled")
	} else {
		if err := OS.DisableStartOnBoot(); err != nil {
			Log.WriteToLogFile("Failed to disable OS-level autostart: " + err.Error())
			return
		}
		Log.WriteToLogFile("Autostart disabled")
	}

	saveAppConfig() // make changes persistent

	Log.WriteToLogFile("StepKeys start on boot state was changed: " + strconv.FormatBool(appConfig.StartOnBoot))

	appConfigMu.Unlock()

	// Notify possible event subscribers
	select {
	case startOnBootChanged <- startOnBootSnapshot:
	default:
	}
}

// Returns the port for the web server
// During normal operation, the web port should not be changed by the user
func GetWebPort() int {
	appConfigMu.RLock()
	defer appConfigMu.RUnlock()
	return appConfig.WebPort
}

// Sets the pedal map
// Used by the API to update the pedal configuration
func SetPedalMap(newMap PedalMap) {
	pedalMapMu.Lock()
	defer pedalMapMu.Unlock()
	pedalMap = newMap

	// Update the pedal map copy in the handler package
	Handler.UpdatePedalMap(newMap)

	data, err := json.MarshalIndent(pedalMap, "", "  ")
	if err != nil {
		Log.WriteToLogFile("Failed to encode pedal map: " + err.Error())
		return
	}

	if err := os.WriteFile(pedalConfigFilePath, data, 0644); err != nil {
		Log.WriteToLogFile("Failed to save pedal config: " + err.Error())
	}

	Log.WriteToLogFile("Pedal map updated and saved.")
}

// Returns a copy of the pedal map
func GetPedalMap() PedalMap {
	pedalMapMu.RLock()
	defer pedalMapMu.RUnlock()

	// Return a copy
	copyMap := make(PedalMap)
	maps.Copy(copyMap, pedalMap)

	return copyMap
}
