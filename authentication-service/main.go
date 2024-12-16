package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/coffeemakingtoaster/water-bottler/authentication-service/pkg/singleton"
	"github.com/coffeemakingtoaster/water-bottler/authentication-service/pkg/utils"
	log "github.com/rs/zerolog/log"
)

var (
	dataBasePath string = "./db.yaml" // Hardcoded path to the db.yaml file
	db           *singleton.DataBaseSingleton
)

// get request to /health
func getHealth(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("/health called")
	fmt.Fprintf(w, "ok")
}

// post request to /checkKey with api key in body
func checkKey(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("/checkKey called")

	// Read the request body
	r_body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Debug().Msg("Request body read")

	// Convert the request body to a string
	// This is the api key
	// Check if the body is not empty and the api key is not over 100 characters long
	api_key := string(r_body)
	if api_key == "" || len(api_key) == 0 || len(api_key) > 100 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check if the api key only contains base64 characters
	if !utils.IsBase64(api_key) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check if the api key is in the database
	for _, key := range db.ApiKeys {
		if key.Key == api_key {
			// Check if the key is not expired
			valid, err := utils.DateInFuture(key.ValidUntil)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			log.Info().Msgf("Api key for %v is valid: %v", key.Name, valid)
			if valid {
				fmt.Fprint(w, "valid")
			} else {
				fmt.Fprint(w, "invalid")
			}
			return
		}
	}
	log.Info().Msg("Key not found")
	fmt.Fprint(w, "invalid")
}

func main() {
	// Get the database instance
	db = singleton.GetDatabaseInstance(dataBasePath)

	// Set up the http server
	http.HandleFunc("/health", getHealth)
	http.HandleFunc("/checkKey", checkKey)

	// Start the server
	interfaceIP := "0.0.0.0"
	interfacePort := 8080
	addr := fmt.Sprintf("%s:%d", interfaceIP, interfacePort)
	log.Info().Msgf("Starting authentication service on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Err(err).Msg("Error starting server")
	}
}
