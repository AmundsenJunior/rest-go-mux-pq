package main

import (
	"io"
	"net/http"
	"github.com/gorilla/handlers"
)

// TODO: create logging formatter that can apply to both App handler and Go main process logs
func (a *App) createLoggingRouter(out io.Writer) http.Handler {
	return handlers.LoggingHandler(out, a.Router)
}
