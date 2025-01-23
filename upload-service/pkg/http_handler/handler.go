package httphandler

import (
	"io"
	"net/http"

	imagestore "github.com/coffeemakingtoaster/water-bottler/upload-service/pkg/image_store"
	"github.com/rs/zerolog/log"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Got health check request")
	ok := true
	if !authServiceIsReachable {
		log.Error().Msg("Auth Service Connection not ok")
	}
	if !imagestore.IsHealthy() {
		log.Error().Msg("Object store connection not available")
	}
	if !ok {
		http.Error(w, "Not ok", http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "ok")
}
