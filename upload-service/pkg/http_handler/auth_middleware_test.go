package httphandler

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

const USERMAIL = "user@mail.com"
const SUCCESS_RESPONSE = `{"status":"valid","email":"valid@solid.com"}`
const INVALID_REPONSE = `{"status":"invalid","email":""}`

func getSampleRequestWithApiKey(key string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Add("X-API-KEY", key)
	return r
}

func getAvailablePort(startPort int) int {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", startPort))

	defer ln.Close()

	if err != nil {
		return getAvailablePort(startPort + 1)
	}

	return startPort
}

func startMockAuthService(validKey string, port int) (*sync.WaitGroup, *http.Server) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/checkKey",
		func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			fmt.Printf("Want: %s Got: %s\n", validKey, string(body))
			if validKey == string(body) {
				fmt.Print(SUCCESS_RESPONSE)
				w.Write([]byte(SUCCESS_RESPONSE))
				return
			}
			w.Write([]byte(INVALID_REPONSE))
		})
	srv := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: serveMux}
	AuthApiUrl = fmt.Sprintf("http://localhost:%d", port)
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

func TestApiKey(t *testing.T) {
	validKey := "valid"
	httpServerExitDone, srv := startMockAuthService(validKey, getAvailablePort(8080))
	w := httptest.NewRecorder()

	calledCount := 0

	protectedHandler := ProtectWithApiKey(func(w http.ResponseWriter, r *http.Request) {
		calledCount++
	})
	protectedHandler(w, getSampleRequestWithApiKey("invalidKey"))

	actual := w.Result()

	checkStatusCodeAndCalledCount(t, actual.StatusCode, http.StatusUnauthorized, calledCount, 0)

	w = httptest.NewRecorder()
	calledCount = 0

	protectedHandler(w, getSampleRequestWithApiKey(validKey))

	actual = w.Result()

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
