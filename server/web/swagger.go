package web

import (
	"net/http"

	// Triggers swagger docs init
	_ "stepkeys/server/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func ServeApiDocs() {
	http.Handle("/api/docs/", httpSwagger.WrapHandler)
}
