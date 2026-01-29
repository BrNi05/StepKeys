package handler

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	"go.bug.st/serial"

	Log "stepkeys/server/logging"
	Pedal "stepkeys/server/pedal"
)

// Pedal state map: pedalID -> pressed: true, released: false
// Only has effect for toggle and hold (behaviour) pedals
var pedalState = make(map[int]bool)

// Currently pressed keyboard keys map: key -> pressed/released
// Used to avoid repeated key down events and to reset state on disable or map change
var keysDown = make(map[string]bool)

// Local copy of the pedal map and enabled state
var (
	pedalMap   = make(Pedal.PedalMap)
	pedalMapMu sync.RWMutex

	enabled   bool
	enabledMu sync.RWMutex
)

// Read the local pedal map that is kept in sync with config pedal map
func readPedalMap() Pedal.PedalMap {
	pedalMapMu.RLock()
	defer pedalMapMu.RUnlock()

	// No copy here
	return pedalMap
}

// Sync the local pedal map with the config pedal map
func UpdatePedalMap(newMap Pedal.PedalMap) {
	pedalMapMu.Lock()
	defer pedalMapMu.Unlock()

	pedalMap = newMap

	// Reset pedals to avoid stuck keys and inconsistent state
	resetPedals()
}

// Read the local enabled state
func readEnabled() bool {
	enabledMu.RLock()
	defer enabledMu.RUnlock()

	return enabled
}

// Sync the local enabled state with the config enabled state
func UpdateEnabled(state bool) {
	enabledMu.Lock()
	defer enabledMu.Unlock()

	enabled = state

	// Reset pedals to avoid stuck keys and inconsistent state
	resetPedals()
}

// Reset all pedals and release any pressed keys
// When StepKeys is disabled or the pedal map is changed this should be called
func resetPedals() {
	for key, pressed := range keysDown {
		if pressed {
			robotgo.KeyUp(key)
			keysDown[key] = false
		}
	}

	for pedalID := range pedalState {
		pedalState[pedalID] = false
	}
}

// Helper: convert []string to []interface{} (array of any) for robotgo
func stringToAny(s []string) []any {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}

// Press and release the keys
func tapKeys(keys []string) {
	for _, key := range keys {
		if !keysDown[key] {
			robotgo.KeyTap(key)
		}
	}
}

// Press the keys down and do not release them
func pressKeys(keys []string) {
	for _, key := range keys {
		if !keysDown[key] {
			robotgo.KeyDown(key)
			keysDown[key] = true
		}
	}
}

// Release the keys
func releaseKeys(keys []string) {
	for _, key := range keys {
		if keysDown[key] {
			robotgo.KeyUp(key)
			keysDown[key] = false
		}
	}
}

// Handle a raw pedal byte from Arduino
func handlePedalByte(b byte) {
	pedalID := int(b & 0x7F)   // Lower 7 bits: pedal ID (0-127)
	pressed := (b & 0x80) != 0 // MSB: pressed/released

	// Only handle pedal IDs defined in the config
	action, ok := readPedalMap()[fmt.Sprintf("%d", pedalID)]
	if !ok {
		Log.WriteToLogFile(fmt.Sprintf("Received unknown pedal ID: %d", pedalID))
		return
	} else {
		Log.WriteToLogFile(fmt.Sprintf("Pedal %d %s", pedalID, map[bool]string{true: "pressed", false: "released"}[pressed]))
	}

	switch action.Behaviour {
	case Pedal.Oneshot:
		// Press event
		if pressed {
			triggerKeys(action)
		}

		// Release event does nothing in oneshot mode

	case Pedal.Toggle:
		if pressed {
			if pedalState[pedalID] {
				// Pressed -> released
				releaseKeys(action.Keys)
				pedalState[pedalID] = false
			} else {
				// Released → pressed
				pressKeys(action.Keys)
				pedalState[pedalID] = true
			}
		}

		// Release event does nothing in toggle mode

	case Pedal.Hold:
		if pressed {
			pressKeys(action.Keys)
			pedalState[pedalID] = true
		} else {
			releaseKeys(action.Keys)
			pedalState[pedalID] = false
		}
	}
}

// Oneshot behaviour helper
func triggerKeys(action Pedal.PedalAction) {
	switch action.Mode {
	case Pedal.Sequence:
		tapKeys(action.Keys)
	case Pedal.Combo:
		if len(action.Keys) == 0 {
			return
		}
		mainKey := action.Keys[len(action.Keys)-1]
		mods := action.Keys[:len(action.Keys)-1]
		robotgo.KeyTap(mainKey, stringToAny(mods)...)
	}
}

// Open the serial port specified in the .env file
// Returns nil if unable to open port
func openSerialPort(baudRate int, serialPort string) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
	}

	port, err := serial.Open(serialPort, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port %s: %w", serialPort, err)
	}

	Log.WriteToLogFile("Serial port opened: " + serialPort)
	return port, nil
}

// Listen to the specified serial port for pedal events
// This is called from main()
func ListenSerial(baudRate int, serialPort string) error {
	port, err := openSerialPort(baudRate, serialPort)
	if err != nil {
		return err
	}
	defer port.Close()

	// Buffer for reading one single byte from the Arduino
	buf := make([]byte, 1)

	// Clear input buffer on start
	port.ResetInputBuffer()

	for {
		// Do not process events if disabled
		// 500ms sleep to avoid high CPU usage
		if !readEnabled() {
			port.ResetInputBuffer()
			time.Sleep(500 * time.Millisecond)
			continue
		}

		n, err := port.Read(buf)
		if err != nil {
			// Ignore errors like timeout or corrupted data
			port.ResetInputBuffer()
			time.Sleep(5 * time.Millisecond)
			continue
		}
		if n == 0 {
			port.ResetInputBuffer() // just to be sure
			continue
		}

		// Check again if StepKeys is enabled
		// This prevents processing of the pedal press that happens after disabling
		if readEnabled() {
			handlePedalByte(buf[0])
		}
	}
}
