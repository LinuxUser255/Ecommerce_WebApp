package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes function having the app Application receiver.
// This Returns an HTTP handler.
// Then create a new router named, "mux", for multiplexing.
// Multiplexing is sending multiple concurrent streams of data.
func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/virtual-terminal", app.VirtualTerminal)

	return mux
}
