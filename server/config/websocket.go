package config

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	Log "stepkeys/server/logging"
)

var (
	settingsclients   = make(map[*websocket.Conn]bool)
	settingsClientsMu sync.Mutex

	pedalsClients   = make(map[*websocket.Conn]bool)
	pedalsClientsMu sync.Mutex

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// Called when enabled or startOnBoot state (may) have changed
func BroadcastSetting(event string, value bool) {
	settingsClientsMu.Lock()
	defer settingsClientsMu.Unlock()

	msg := fmt.Sprintf(`{"event":"%s","value":%t}`, event, value)
	for client := range settingsclients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			client.Close()
			delete(settingsclients, client)
		}
	}
}

// WebSocket handler for settings
func SettingsWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Log.WriteToLogFile("Settings WebSocket upgrade error: " + err.Error())
	}

	settingsClientsMu.Lock()
	settingsclients[conn] = true
	settingsClientsMu.Unlock()

	// Keep the connection alive
	for {
		if _, _, err := conn.NextReader(); err != nil {
			settingsClientsMu.Lock()
			delete(settingsclients, conn)
			settingsClientsMu.Unlock()
			conn.Close()
			break
		}
	}
}

// Called when the pedal map is updated
func NotifyPedalMapUpdate() {
	pedalsClientsMu.Lock()
	defer pedalsClientsMu.Unlock()

	for client := range pedalsClients {
		// Message doesn't matter, this acts as a change ping
		if err := client.WriteMessage(websocket.TextMessage, []byte("1")); err != nil {
			client.Close()
			delete(pedalsClients, client)
		}
	}
}

// WebSocket handler for pedal map updates
func PedalWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Log.WriteToLogFile("Pedals WebSocket upgrade error: " + err.Error())
		return
	}

	pedalsClientsMu.Lock()
	pedalsClients[conn] = true
	pedalsClientsMu.Unlock()

	// Keep the connection alive
	for {
		if _, _, err := conn.NextReader(); err != nil {
			pedalsClientsMu.Lock()
			delete(pedalsClients, conn)
			pedalsClientsMu.Unlock()
			conn.Close()
			break
		}
	}
}
