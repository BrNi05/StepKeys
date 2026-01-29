package web

import (
	"net/http"

	// Triggers swagger docs init
	_ "stepkeys/server/docs"
	Log "stepkeys/server/logging"

	httpSwagger "github.com/swaggo/http-swagger"
)

func ServeApiDocs() {
	http.Handle("/api/docs/", httpSwagger.WrapHandler)

	Log.WriteToLogFile("API docs route registered at /api/docs/.")
}
