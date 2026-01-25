package web

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	Config "stepkeys/server/config"
	. "stepkeys/server/pedal"
	Updater "stepkeys/server/updater"
)

const (
	methodNotAllowed = "Method not allowed."
	contentType      = "Content-Type"
	contentTypeJson  = "application/json"
)

// @Description Standard JSON error response
type ErrorResponse struct {
	Status       int    `json:"status" example:"400"`
	ErrorMessage string `json:"errorMessage" example:"Invalid JSON payload"`
}

// @Description Standard JSON boolean response
type BooleanResponse struct {
	Value bool `json:"value" example:"true"`
}

// Helper: construct and send JSON error response
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Status:       status,
		ErrorMessage: msg,
	})
}

// @Summary      Get all pedals
// @Description  Returns the full pedal configuration map.
// @Tags         pedals
// @Produce      json
// @Success      200 {object} PedalMap
// @Router       /api/pedals [get]
func getPedals(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(contentType, contentTypeJson)
	_ = json.NewEncoder(w).Encode(Config.GetPedalMap())
}

// @Summary      Update all pedals
// @Description  Replaces the current pedal configuration entirely.
// @Tags         pedals
// @Accept       json
// @Produce      json
// @Param        pedals  body  PedalMap  true  "New pedal configuration"
// @Success      200     {object} PedalMap
// @Failure      400     {object} ErrorResponse
// @Router       /api/pedals [post]
func updatePedals(w http.ResponseWriter, r *http.Request) {
	var newConfig PedalMap

	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validate the new config
	if err := ValidatePedalMap(newConfig); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid pedal configuration: "+err.Error())
		return
	}

	Config.SetPedalMap(newConfig)

	w.Header().Set(contentType, contentTypeJson)
	_ = json.NewEncoder(w).Encode(newConfig)
}

// @Summary      Get enabled state
// @Description  Returns whether StepKeys is enabled or disabled.
// @Tags         settings
// @Produce      json
// @Success      200 {object} BooleanResponse
// @Router       /api/enabled [get]
func getEnabled(w http.ResponseWriter, _ *http.Request) {
	enabled := Config.IsEnabled()
	w.Header().Set(contentType, contentTypeJson)
	_ = json.NewEncoder(w).Encode(BooleanResponse{Value: enabled})
}

// @Summary      Toggle enabled state
// @Description  Toggles whether StepKeys is enabled or disabled. If the pedal configuration is empty, enabling will silently fail.
// @Tags         settings
// @Produce      json
// @Success      200 {object} BooleanResponse
// @Router       /api/enabled [post]
func toggleEnabled(w http.ResponseWriter, _ *http.Request) {
	Config.ToggleEnabled()
	w.Header().Set(contentType, contentTypeJson)
	_ = json.NewEncoder(w).Encode(BooleanResponse{Value: Config.IsEnabled()})
}

// @Summary      Get start on boot state
// @Description  Returns whether StepKeys is set to start on boot.
// @Tags         settings
// @Produce      json
// @Success      200 {object} BooleanResponse
// @Router       /api/boot [get]
func getStartOnBootEnabled(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(contentType, contentTypeJson)
	_ = json.NewEncoder(w).Encode(BooleanResponse{Value: Config.IsStartOnBootEnabled()})
}

// @Summary      Toggle start on boot state
// @Description  Toggles whether StepKeys should start on boot.
// @Tags         settings
// @Produce      json
// @Success      200 {object} BooleanResponse
// @Router       /api/boot [post]
func toggleStartOnBoot(w http.ResponseWriter, _ *http.Request) {
	Config.ToggleStartOnBoot()
	w.Header().Set(contentType, contentTypeJson)
	_ = json.NewEncoder(w).Encode(BooleanResponse{Value: Config.IsStartOnBootEnabled()})
}

// @Summary      Quit application
// @Description  Terminates the StepKeys process.
// @Tags         lifecycle
// @Success      200
// @Router       /api/quit [post]
func quitApp(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()
}

// @Summary      Check for updates
// @Description  Returns whether a new update is available.
// @Tags         lifecycle
// @Produce      json
// @Param        force query bool false "Force check for updates"
// @Success      200 {object} BooleanResponse
// @Router       /api/update [get]
func getUpdateAvailable(w http.ResponseWriter, r *http.Request) {
	forceCheck := r.URL.Query().Get("force")

	if forceCheck == "1" || forceCheck == "true" {
		Updater.CheckForUpdates()
	}

	w.Header().Set(contentType, contentTypeJson)
	_ = json.NewEncoder(w).Encode(BooleanResponse{Value: Updater.UpdateAvailable()})
}

// Registers all API routes
func RegisterAPI() {
	http.HandleFunc("/api/pedals", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getPedals(w, r)
		case http.MethodPost:
			updatePedals(w, r)
		default:
			writeJSONError(w, http.StatusMethodNotAllowed, methodNotAllowed)
		}
	})

	http.HandleFunc("/api/enabled", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getEnabled(w, r)
		case http.MethodPost:
			toggleEnabled(w, r)
		default:
			writeJSONError(w, http.StatusMethodNotAllowed, methodNotAllowed)
		}
	})

	http.HandleFunc("/api/boot", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getStartOnBootEnabled(w, r)
		case http.MethodPost:
			toggleStartOnBoot(w, r)
		default:
			writeJSONError(w, http.StatusMethodNotAllowed, methodNotAllowed)
		}
	})

	http.HandleFunc("/api/quit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONError(w, http.StatusMethodNotAllowed, methodNotAllowed)
			return
		}
		quitApp(w, r)
	})

	http.HandleFunc("/api/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSONError(w, http.StatusMethodNotAllowed, methodNotAllowed)
			return
		}
		getUpdateAvailable(w, r)
	})

	log.Println("API routes registered.")
}
