package httphandler

import (
	"io"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Got health check request")
	// Check if the service has access to the object store
	if !objectStoreConnAvailable {
		log.Error().Msg("Object store connection not available")
		http.Error(w, "Not ok", http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "ok")
}
