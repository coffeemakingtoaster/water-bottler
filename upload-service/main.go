package main

import (
	"fmt"
	"net/http"

	httphandler "github.com/coffeemakingtoaster/water-bottler/upload-service/pkg/http_handler"
)

func main() {
	http.HandleFunc("/health", httphandler.GetHealth)
	http.HandleFunc("/upload", httphandler.ProtectWithApiKey(httphandler.HandleUpload))
	fmt.Println("Starting upload service")
	err := http.ListenAndServe(":8081", nil)
	fmt.Printf("Server encountered error: %v", err)
}
