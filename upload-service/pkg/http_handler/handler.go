package httphandler

import (
	"io"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}
