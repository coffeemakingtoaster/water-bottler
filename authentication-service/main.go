package main

import (
	"fmt"
	"io"
	"net/http"
)

func getHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got health request")
	io.WriteString(w, "ok")
}

func main() {
	http.HandleFunc("/health", getHealth)
	fmt.Println("Starting authentication service")
	err := http.ListenAndServe(":8080", nil)
	fmt.Printf("Server encountered error: %v", err)
}
