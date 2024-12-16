package httphandler

import (
	"errors"
	"net/http"
	"time"
)

type CachedApiKey struct {
	Email     string
	FetchedAt time.Time
}

type ApiKeyCache struct {
	cache map[string]CachedApiKey
}

// How long can an api key be cached before we throw it out
const API_KEY_TTL = time.Minute * 2

const USER_EMAIL_HEADER = "X-user-mail"

var apiKeyCache = ApiKeyCache{map[string]CachedApiKey{}}

func (akc *ApiKeyCache) retrieveValidKey(key string) (string, error) {
	val, ok := akc.cache[key]
	if !ok {
		return "", errors.New("Api Key not present")
	}
	if time.Time.Sub(time.Now(), val.FetchedAt) > API_KEY_TTL {
		return "", errors.New("Cached entry expired")
	}
	return val.Email, nil
}

func (akc *ApiKeyCache) addKeyToCache(key, mail string) {
	akc.cache[key] = CachedApiKey{
		mail,
		time.Now(),
	}
}

func validateAPIKeyViaAuthService(key string) (string, error) {
	if key == "aaaaa" {
		return "test@test.de", nil
	}
	return "", errors.New("Invalid API key")
}

func hasValidApiKey(r *http.Request) (bool, string) {
	apiKey := r.Header.Get("X-API-KEY")

	// Do we have a valid entry cached?
	mail, err := apiKeyCache.retrieveValidKey(apiKey)
	if err == nil {
		return true, mail
	}

	mail, err = validateAPIKeyViaAuthService(apiKey)
	if err != nil {
		return false, ""
	}

	apiKeyCache.addKeyToCache(apiKey, mail)

	return true, mail
}

func ProtectWithApiKey(handler HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, mail := hasValidApiKey(r)

		// save for later stages
		r.Header.Add(USER_EMAIL_HEADER, mail)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}
