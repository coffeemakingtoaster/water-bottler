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

func startMockAuthService(expectedResponse string) (*sync.WaitGroup, *http.Server) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

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

	return wg, &srv
}

func stopMockAuthService(wg *sync.WaitGroup, srv *http.Server) {
	srv.Close()
	wg.Wait()
}

func checkStatusCodeAndCalledCount(t *testing.T, actualStatusCode, expectedStatusCode, actualCalledCount, expectedCalledCount int) {
	if actualStatusCode != expectedStatusCode {
		t.Errorf("Expected Response of %d but got %d", expectedStatusCode, actualStatusCode)
	}

	if actualCalledCount != expectedCalledCount {
		t.Errorf("Expected protected handler to be called %d times but was called %d times", expectedCalledCount, actualCalledCount)
	}
}

func TestCachedApiKey(t *testing.T) {
	apiKey := "cachedKey123"
	apiKeyCache.addKeyToCache(apiKey, USERMAIL)
	w := httptest.NewRecorder()

	actualMail := ""

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		actualMail = r.Header.Get(USER_EMAIL_HEADER)
	})
	protectedHandler(w, getSampleRequestWithApiKey(apiKey))

	actual := w.Result()

	if actual.StatusCode != http.StatusOK {
		t.Errorf("Expected Response of %d but got %d", http.StatusOK, actual.StatusCode)
	}

	if actualMail != USERMAIL {
		t.Errorf("Expected mail in header of %s but got %s", USERMAIL, actualMail)
	}
}

func TestInvalidApiKey(t *testing.T) {
	httpServerExitDone, srv := startMockAuthService(`{"status":"invalid","email":""}`)
	w := httptest.NewRecorder()

	calledCount := 0

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		calledCount++
	})
	protectedHandler(w, getSampleRequestWithApiKey("invalidKey"))

	actual := w.Result()

	checkStatusCodeAndCalledCount(t, actual.StatusCode, http.StatusUnauthorized, calledCount, 0)

	stopMockAuthService(httpServerExitDone, srv)
}

func TestValidApiKey(t *testing.T) {
	httpServerExitDone, srv := startMockAuthService(`{"status":"valid","email":"valid@solid.com"}`)

	w := httptest.NewRecorder()

	calledCount := 0

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		calledCount++
	})
	protectedHandler(w, getSampleRequestWithApiKey("valid"))

	actual := w.Result()

	checkStatusCodeAndCalledCount(t, actual.StatusCode, http.StatusOK, calledCount, 1)

	stopMockAuthService(httpServerExitDone, srv)
}

func TestNoApiKey(t *testing.T) {
	w := httptest.NewRecorder()

	calledCount := 0

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		calledCount++
	})
	protectedHandler(w, httptest.NewRequest(http.MethodGet, "/", nil))

	actual := w.Result()

	checkStatusCodeAndCalledCount(t, actual.StatusCode, http.StatusUnauthorized, calledCount, 0)
}

func TestUnreachableAuthService(t *testing.T) {
	w := httptest.NewRecorder()

	calledCount := 0

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		calledCount++
	})
	protectedHandler(w, getSampleRequestWithApiKey("doesntmatter"))

	actual := w.Result()
	checkStatusCodeAndCalledCount(t, actual.StatusCode, http.StatusInternalServerError, calledCount, 0)
}
