package web

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	Config "stepkeys/server/config"
)

//go:embed gui/*
var guiFS embed.FS

func StartGUI() {
	// Sub filesystem to serve UI at root
	subFS, err := fs.Sub(guiFS, "gui")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(subFS)))

	addr := fmt.Sprintf(":%d", Config.GetWebPort())
	log.Println("webGUI started on port", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("Failed to start webGUI:", err)
	}
}
