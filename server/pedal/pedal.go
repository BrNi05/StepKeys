package pedal

import (
	"fmt"
)

// Pedal mode
type PedalMode string

const (
	Sequence PedalMode = "sequence"
	Combo    PedalMode = "combo"
)

// Pedal behaviour
type PedalBehaviour string

const (
	Oneshot PedalBehaviour = "oneshot"
	Toggle  PedalBehaviour = "toggle"
)

// PedalAction describes what a pedal does when pressed
type PedalAction struct {
	// Mode defines how keys are triggered
	// sequence = keys pressed one after another
	// combo = keys pressed together (a key combo)
	Mode PedalMode `json:"mode" example:"sequence"`

	// Keys are the key names sent to the OS
	// Supported: https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys
	Keys []string `json:"keys" example:"ctrl,shift,escape"`

	// Behaviour defines how the pedal behaves while pressed
	// oneshot = press keys once per pedal press
	// toggle = the keys are held down until the pedal is pressed again
	Behaviour PedalBehaviour `json:"behaviour" example:"oneshot"`
}

// PedalMap represents the full pedal configuration.
// @Description Map of pedal IDs to their assigned actions
// @example {"pedal_1":{"mode":"sequence","keys":["ctrl","shift"],"behaviour":"oneshot"}}
type PedalMap map[string]PedalAction

// Validate the pedal mode string
func isValidMode(mode PedalMode) bool {
	return mode == Sequence || mode == Combo
}

// Validate the pedal behaviour string
func isValidBehaviour(behaviour PedalBehaviour) bool {
	return behaviour == Oneshot || behaviour == Toggle
}

// Checks if all keys are valid
// Supported: https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys
func isValidKeys(keys []string) bool {
	for _, key := range keys {
		// Single-character keys: A-Z, a-z, 0-9
		if len(key) == 1 {
			r := rune(key[0])
			if (r >= 'a' && r <= 'z') ||
				(r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') {
				continue
			}
			return false
		}

		if _, ok := validNamedKeys[key]; !ok {
			return false
		}
	}
	return true
}

func ValidatePedalMap(m PedalMap) error {
	for pedalID, action := range m {
		if !isValidMode(action.Mode) {
			return fmt.Errorf("Pedal %q: invalid mode %q (use <sequence> or <combo>)",
				pedalID, action.Mode)
		}
		if !isValidBehaviour(action.Behaviour) {
			return fmt.Errorf("Pedal %q: invalid behaviour %q (use <oneshot> or <toggle>)",
				pedalID, action.Behaviour)
		}

		if !isValidKeys(action.Keys) {
			return fmt.Errorf("Pedal %q: contains invalid keys", pedalID)
		}
	}
	return nil
}
