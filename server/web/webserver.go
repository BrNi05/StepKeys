package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	Config "stepkeys/server/config"
	Log "stepkeys/server/logging"
)

//go:embed gui/*
var guiFS embed.FS

func StartGUI() {
	// Sub filesystem to serve UI at root
	subFS, err := fs.Sub(guiFS, "gui")
	if err != nil {
		Log.ErrorToLogFile(fmt.Sprintf("Failed to start webGUI: %v", err))
	}

	http.Handle("/", http.FileServer(http.FS(subFS)))

	addr := fmt.Sprintf(":%d", Config.GetWebPort())
	Log.WriteToLogFile(fmt.Sprintf("webGUI started on port %s.", addr))
	if err := http.ListenAndServe(addr, nil); err != nil {
		Log.ErrorToLogFile(fmt.Sprintf("Failed to start webGUI: %v", err))
	}
}
