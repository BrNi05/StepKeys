package updater

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Set to true if an update is available, false otherwise or if check is disabled
var (
	updateAvailable   bool
	updateAvailableMu sync.RWMutex
)

// Fetches the latest release and compares it with the current version
func CheckForUpdates() {
	updateAvailableMu.Lock()
	defer updateAvailableMu.Unlock()

	// Skip version check if disabled
	if os.Getenv("NO_VERSION_CHECK") != "" {
		log.Println("Version check is disabled.")
		updateAvailable = false
		return
	}

	const repoLatest = "https://github.com/BrNi05/StepKeys/releases/latest"

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Head(repoLatest)
	if err != nil {
		updateAvailable = false
		return
	}
	defer response.Body.Close()

	location := response.Header.Get("Location")
	if location == "" {
		updateAvailable = false
		return
	}

	parts := strings.Split(location, "/")
	latestVersion := strings.TrimPrefix(parts[len(parts)-1], "v")

	// Current version
	currentVersion := os.Getenv("VERSION")

	// Compare versions (lexicographically)
	if currentVersion == "" || latestVersion > currentVersion {
		updateAvailable = true
	} else {
		updateAvailable = false
	}

	log.Println("Update check complete. Current version:", currentVersion, "Latest version:", latestVersion, "Update available:", updateAvailable)
}

func UpdateAvailable() bool {
	updateAvailableMu.RLock()
	defer updateAvailableMu.RUnlock()

	return updateAvailable
}
