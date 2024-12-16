package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// This function is used to test the http handler
func TestHttpHandler(t *testing.T, handler http.HandlerFunc, requestMethod string, body io.Reader, wantedCode int, wantedBody string) {
	req := httptest.NewRequest(requestMethod, "/checkKey", body)
	rec := httptest.NewRecorder()
	handler(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != wantedCode {
		t.Errorf("checkKey() = %v; want %v", res.StatusCode, wantedCode)
	}
	if body, _ := io.ReadAll(res.Body); string(body) != wantedBody {
		t.Errorf("checkKey() = %v; want %v", string(body), wantedBody)
	}
}
