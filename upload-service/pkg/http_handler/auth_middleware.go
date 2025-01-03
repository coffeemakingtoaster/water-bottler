package httphandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	customerror "github.com/coffeemakingtoaster/water-bottler/upload-service/pkg/custom_error"
)

type CachedApiKey struct {
	Email     string
	FetchedAt time.Time
}

type ApiKeyCache struct {
	cache map[string]CachedApiKey
}

type KeyCheckResponse struct {
	Status string `json:"status"`
	Email  string `json:"email"`
}

var AuthApiUrl string

var httpClient = http.Client{}

// How long can an api key be cached before we throw it out
const API_KEY_TTL = time.Minute * 2

const USER_EMAIL_HEADER = "X-user-mail"

const AUTH_SERVICE_URL_ENV = "AUTH_SERVICE_URL"

var apiKeyCache = ApiKeyCache{map[string]CachedApiKey{}}

func init() {
	AuthApiUrl = os.Getenv(AUTH_SERVICE_URL_ENV)
	if len(AuthApiUrl) == 0 {
		AuthApiUrl = "http://localhost:8080"
	}
}

func (akc *ApiKeyCache) retrieveValidKey(key string) (string, error) {
	val, ok := akc.cache[key]
	if !ok {
		return "", customerror.NewSafeErrorFromError(errors.New("Api Key not present"))
	}
	if time.Time.Sub(time.Now(), val.FetchedAt) > API_KEY_TTL {
		return "", customerror.NewSafeErrorFromError(errors.New("Api Key cache entry expired"))
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

	requestUrl := fmt.Sprintf("%s/checkKey", AuthApiUrl)

	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(key))
	if err != nil {
		fmt.Printf("Could init request to auth service with url %s %v\n", requestUrl, err)
		return "", &customerror.SystemCommunicationError{Reason: "Error creating request"}
	}

	res, err := httpClient.Do(req)

	if res == nil || err != nil {
		fmt.Println("Could not complete request")
		return "", &customerror.SystemCommunicationError{Reason: "Could not comlete request"}
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Got invalid response code %d\n", res.StatusCode)
		return "", &customerror.SystemCommunicationError{Reason: fmt.Sprintf("Got invalid response code %d", res.StatusCode)}
	}

	var resp KeyCheckResponse

	json.NewDecoder(res.Body).Decode(&resp)

	if resp.Status != "valid" {
		return "", customerror.NewSafeErrorFromError(errors.New("Invalid API Key"))
	}

	return resp.Email, nil
}

func getEmailForApiKey(r *http.Request) (string, error) {
	apiKey := r.Header.Get("X-API-KEY")

	if len(apiKey) == 0 {
		return "", customerror.NewSafeErrorFromError(errors.New("No API Key provided"))
	}

	// Do we have a valid entry cached?
	mail, err := apiKeyCache.retrieveValidKey(apiKey)
	if err == nil {
		return mail, nil
	}

	mail, err = validateAPIKeyViaAuthService(apiKey)
	if err == nil {
		apiKeyCache.addKeyToCache(apiKey, mail)
		return mail, nil
	}

	return "", err
}

func ProtectWithApiKey(handler HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mail, err := getEmailForApiKey(r)

		if err != nil {
			switch err := err.(type) {
			case *customerror.SystemCommunicationError:
				w.WriteHeader(http.StatusInternalServerError)
			case *customerror.SafeError:
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.OutwardMessage))
			default:
				w.WriteHeader(http.StatusUnauthorized)
			}
			return
		}

		// save for later stages
		r.Header.Add(USER_EMAIL_HEADER, mail)

		handler(w, r)
	}
}
