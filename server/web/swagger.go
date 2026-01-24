package web

import (
	"log"
	"net/http"

	// Triggers swagger docs init
	_ "stepkeys/server/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func ServeApiDocs() {
	http.Handle("/api/docs/", httpSwagger.WrapHandler)

	log.Println("API docs route registered at /api/docs/.")
}
