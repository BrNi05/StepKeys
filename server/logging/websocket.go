package logging

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// Called when a new log message is generated
func broadcastLog(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			// Client is no longer connected
			client.Close()
			delete(clients, client)
		}
	}
}

// Serve WebSocket connections
func LogsWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ErrorToLogFile("Logs WebSocket upgrade error: " + err.Error()) // fatal
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	// Keep connection open
	for {
		if _, _, err := conn.NextReader(); err != nil {
			clientsMu.Lock()
			delete(clients, conn)
			clientsMu.Unlock()
			conn.Close()
			break
		}
	}
}
