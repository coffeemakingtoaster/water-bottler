package httphandler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const USERMAIL = "user@mail.com"

func getSampleRequestWithApiKey(key string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Add("X-API-KEY", key)
	return r
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
