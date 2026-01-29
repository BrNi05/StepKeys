package config

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	Log "stepkeys/server/logging"
)

var (
	settingsClients   = make(map[*websocket.Conn]bool)
	settingsClientsMu sync.Mutex
	upgrader          = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "http://localhost:18000" || origin == "http://127.0.0.1:18000"
		},
	}
)

// Called when enabled or startOnBoot state (may) have changed
func BroadcastSetting(event string, value bool) {
	settingsClientsMu.Lock()
	defer settingsClientsMu.Unlock()

	msg := fmt.Sprintf(`{"event":"%s","value":%t}`, event, value)
	for client := range settingsClients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			client.Close()
			delete(settingsClients, client)
		}
	}
}

// WebSocket handler for settings
func SettingsWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Log.ErrorToLogFile("Settings WebSocket upgrade error: " + err.Error()) // fatal
	}

	settingsClientsMu.Lock()
	settingsClients[conn] = true
	settingsClientsMu.Unlock()

	// Keep connection open
	for {
		if _, _, err := conn.NextReader(); err != nil {
			settingsClientsMu.Lock()
			delete(settingsClients, conn)
			settingsClientsMu.Unlock()
			conn.Close()
			break
		}
	}
}
