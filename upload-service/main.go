package main

import (
	"fmt"
	"net/http"

	httphandler "github.com/coffeemakingtoaster/water-bottler/upload-service/pkg/http_handler"
	"github.com/rs/zerolog/log"
)

func main() {
	http.HandleFunc("/health", httphandler.GetHealth)
	http.HandleFunc("/upload", httphandler.ProtectWithApiKey(httphandler.HandleUpload))

	interfaceIP := "0.0.0.0"
	interfacePort := 8081
	addr := fmt.Sprintf("%s:%d", interfaceIP, interfacePort)

	log.Info().Msgf("Starting upload service on %s", addr)
	err := http.ListenAndServe(addr, nil)
	log.Error().Msgf("Server encountered error: %v", err)
}
