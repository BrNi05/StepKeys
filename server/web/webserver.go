package web

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	Config "stepkeys/server/config"
	. "stepkeys/server/os"
)

// Serve the webGUI
func StartGUI() {
	webDir := filepath.Join(filepath.Dir(GetExecPath()), "webgui")

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	addr := fmt.Sprintf(":%d", Config.GetWebPort())
	log.Println("Starting webGUI on port", addr, "serving folder:", webDir)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("Failed to start webGUI:", err)
	}
}
