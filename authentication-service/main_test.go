package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/coffeemakingtoaster/water-bottler/authentication-service/pkg/models"
	"github.com/coffeemakingtoaster/water-bottler/authentication-service/pkg/singleton"
	"github.com/coffeemakingtoaster/water-bottler/authentication-service/pkg/utils"
)

func TestDateInFuture(t *testing.T) {
	t.Log("Check for date in the future")
	futureDate := time.Now().AddDate(0, 0, +1)
	if inFuture, _ := utils.DateInFuture(futureDate); inFuture != true {
		t.Errorf("DateInFuture() = %v; want true", inFuture)
	}
	t.Log("Check for date in the past")
	pastDate := time.Now().AddDate(0, 0, -1)
	if inFuture, _ := utils.DateInFuture(pastDate); inFuture != false {
		t.Errorf("DateInFuture() = %v; want false", inFuture)
	}
	t.Log("Check for current date")
	nowDate := time.Now()
	if inFuture, _ := utils.DateInFuture(nowDate); inFuture != false {
		t.Errorf("DateInFuture() = %v; want false", inFuture)
	}
}

func TestCheckKey(t *testing.T) {
	endpoint := "/checkKey"
	db = singleton.GetDatabaseInstance("./db.yaml.examble")
	var validKey singleton.ApiKey
	var expiredKey singleton.ApiKey
	for _, key := range db.ApiKeys {
		validUntil, _ := time.Parse(time.RFC3339, key.ValidUntil)
		if validUntil.After(time.Now()) {
			validKey = key
		} else {
			expiredKey = key
		}
		// Break if we have a valid and invalid key
		if validKey.Key != "" && expiredKey.Key != "" {
			break
		}
	}

	if validKey.Key == "" {
		t.Log("No valid key found, skipping test for a valid key")
	} else {
		t.Log("Check for valid key")
		utils.TestHttpHandler(t, checkKey, "POST", endpoint, strings.NewReader(validKey.Key), 200, fmt.Sprintf(`{"status":"valid","email":"%s"}`, validKey.Name))
	}

	if expiredKey.Key == "" {
		t.Log("No expired key found, skipping test for an expired key")
	} else {
		t.Log("Check for invalid key")
		utils.TestHttpHandler(t, checkKey, "POST", endpoint, strings.NewReader(expiredKey.Key), 200, `{"status":"invalid","email":""}`)
	}

	t.Log("Check for non-existing key")
	utils.TestHttpHandler(t, checkKey, "POST", endpoint, strings.NewReader("nonExistingKey"), 200, `{"status":"invalid","email":""}`)

	t.Log("Check for invalid request")
	utils.TestHttpHandler(t, checkKey, "POST", endpoint, models.ErrorReader{}, 500, http.StatusText(500)+"\n")
}
