package httphandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

const USERMAIL = "user@mail.com"

func getSampleRequestWithApiKey(key string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Add("X-API-KEY", key)
	return r
}

func startMockAuthService(wg *sync.WaitGroup, expectedResponse string) *http.Server {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/checkKey",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expectedResponse))

		})
	srv := http.Server{Addr: ":8080", Handler: serveMux}
	AuthApiUrl = "http://localhost:8080"
	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(fmt.Sprintf("ListenAndServe(): %v", err))
		}
	}()

	return &srv
}

func stopMockAuthService(wg *sync.WaitGroup, srv *http.Server) {
	srv.Close()
	wg.Wait()
}

func TestCachedApiKey(t *testing.T) {
	apiKey := "cachedKey123"
	req := getSampleRequestWithApiKey(apiKey)
	apiKeyCache.addKeyToCache(apiKey, USERMAIL)
	w := httptest.NewRecorder()

	actualMail := ""
	expectedResponseCode := http.StatusOK

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		actualMail = r.Header.Get(USER_EMAIL_HEADER)
	})
	protectedHandler(w, req)

	actual := w.Result()

	if actual.StatusCode != expectedResponseCode {
		t.Errorf("Expected Response of %d but got %d", expectedResponseCode, actual.StatusCode)
	}

	if actualMail != USERMAIL {
		t.Errorf("Expected mail in header of %s but got %s", USERMAIL, actualMail)
	}
}

func TestInvalidApiKey(t *testing.T) {
	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)
	srv := startMockAuthService(httpServerExitDone, `{"status":"invalid","email":""}`)
	req := getSampleRequestWithApiKey("invalidKey")
	w := httptest.NewRecorder()

	expectedResponseCode := http.StatusUnauthorized
	calledCount := 0

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		calledCount++
	})
	protectedHandler(w, req)

	actual := w.Result()

	if actual.StatusCode != expectedResponseCode {
		t.Errorf("Expected Response of %d but got %d", expectedResponseCode, actual.StatusCode)
	}

	if calledCount > 0 {
		t.Errorf("Expected protected handler to not be called. However handler was called %d times", calledCount)
	}

	stopMockAuthService(httpServerExitDone, srv)
}

func TestValidApiKey(t *testing.T) {
	req := getSampleRequestWithApiKey("valid")
	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)
	srv := startMockAuthService(httpServerExitDone, `{"status":"valid","email":"valid@solid.com"}`)

	w := httptest.NewRecorder()

	expectedResponseCode := http.StatusOK
	calledCount := 0

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		calledCount++
	})
	protectedHandler(w, req)

	actual := w.Result()

	if actual.StatusCode != expectedResponseCode {
		t.Errorf("Expected Response of %d but got %d", expectedResponseCode, actual.StatusCode)
	}

	if calledCount != 1 {
		t.Errorf("Expected protected handler have been called once. However handler was called %d times", calledCount)
	}
	stopMockAuthService(httpServerExitDone, srv)
}

func TestNoApiKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	expectedResponseCode := http.StatusUnauthorized
	calledCount := 0

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		calledCount++
	})
	protectedHandler(w, req)

	actual := w.Result()

	if actual.StatusCode != expectedResponseCode {
		t.Errorf("Expected Response of %d but got %d", expectedResponseCode, actual.StatusCode)
	}

	if calledCount > 0 {
		t.Errorf("Expected protected handler to not be called. However handler was called %d times", calledCount)
	}
}

func TestUnreachableAuthService(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)
	srv := startMockAuthService(httpServerExitDone, `{"status":"invalid","email":""}`)
	defer stopMockAuthService(httpServerExitDone, srv)

	w := httptest.NewRecorder()

	expectedResponseCode := http.StatusUnauthorized
	calledCount := 0

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		calledCount++
	})
	protectedHandler(w, req)

	actual := w.Result()

	if actual.StatusCode != expectedResponseCode {
		t.Errorf("Expected Response of %d but got %d", expectedResponseCode, actual.StatusCode)
	}

	if calledCount > 0 {
		t.Errorf("Expected protected handler to not be called. However handler was called %d times", calledCount)
	}
}
