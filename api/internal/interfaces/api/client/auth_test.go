package client_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/buffer"
	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/server"
	"github.com/kodmain/thetiptop/api/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

type TokenStructure struct {
	JWT string `json:"jwt"`
}

const (
	GOOD_EMAIL = "user1@example.com"
	GOOD_PASS  = "ValidP@ssw0rd1"

	WRONG_EMAIL = "user2@example.com"
	WRONG_PASS  = "secret"
)

var srv *server.Server

func start() error {
	config.Load("../../../../config.test.yml")
	logger.Info("starting application")
	srv = server.Create()
	srv.Register(interfaces.Endpoints)
	return srv.Start()
}

func stop() error {
	logger.Info("waiting for application to shutdown")
	return srv.Stop()
}

func TestSignUp(t *testing.T) {
	assert.Nil(t, start())

	users := []struct {
		email    string
		password string
		status   int
	}{
		{GOOD_EMAIL, GOOD_PASS, http.StatusCreated},
		{GOOD_EMAIL, GOOD_PASS, http.StatusConflict},
		{WRONG_EMAIL, WRONG_PASS, http.StatusBadRequest},
	}

	for _, user := range users {
		// Create a form with email and password fields
		form := url.Values{}
		form.Set("email", user.email)
		form.Set("password", user.password)

		// Create a new HTTP request to call /sign/up
		req, err := http.NewRequest("POST", "http://localhost/sign/up", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		// Set the form as the request body
		req.Body = io.NopCloser(strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Perform the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to perform request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != user.status {
			t.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}

	assert.Nil(t, stop())
}

func TestSignIn(t *testing.T) {
	TestSignUp(t)
	assert.Nil(t, start())

	users := []struct {
		email    string
		password string
		status   int
	}{
		{GOOD_EMAIL, GOOD_PASS, http.StatusOK},
		{GOOD_EMAIL, WRONG_PASS, http.StatusBadRequest},
		{GOOD_PASS, WRONG_PASS, http.StatusBadRequest},
	}

	for _, user := range users {
		// Create a form with email and password fields
		form := url.Values{}
		form.Set("email", user.email)
		form.Set("password", user.password)

		// Create a new HTTP request to call /sign/in
		req, err := http.NewRequest("POST", "http://localhost/sign/in", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		// Set the form as the request body
		req.Body = io.NopCloser(strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Perform the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to perform request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != user.status {
			t.Errorf("unexpected status code: %d", resp.StatusCode)
		}

	}

	assert.Nil(t, stop())
}

func TestSignRenew(t *testing.T) {
	TestSignUp(t)
	assert.Nil(t, start())

	// Create a new HTTP request to call /sign/in
	req, err := http.NewRequest("POST", "http://localhost/sign/in", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	form := url.Values{}
	form.Set("email", GOOD_EMAIL)
	form.Set("password", GOOD_PASS)

	// Set the form as the request body
	req.Body = io.NopCloser(strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	token, err := buffer.Read(resp.Body)
	assert.Nil(t, err)

	// Déclaration de la variable qui recevra la valeur désérialisée
	var tokenData TokenStructure

	// Désérialisation du JSON dans la structure définie
	err = json.Unmarshal(token.Bytes(), &tokenData)
	if err != nil {
		fmt.Printf("Error while parsing JSON: %s\n", err)
		return
	}

	access, err := serializer.FromString(tokenData.JWT)
	if err != nil {
		t.Error(err)
	}

	users := []struct {
		token  string
		status int
	}{
		{"Bearer " + *access.Refresh, http.StatusOK}, // Replace with actual valid JWT token
		{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDkxMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJ0eXBlIjoxLCJ0eiI6IkxvY2FsIn0.5Lae56HNcQ1OHcP_FhTfcOOtHpaZVgRFy6vzzBugN7Y", http.StatusUnauthorized}, // Replace with actual expired JWT token
		{"Bearer malformed.jwt.token.here", http.StatusUnauthorized}, // Replace with actual malformed JWT token
		{"", http.StatusBadRequest},                                  // Replace with actual empty JWT token
	}

	for _, user := range users {
		// Create a new HTTP request to call /sign/renew
		req, err := http.NewRequest("GET", "http://localhost/sign/renew", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		// Set the JWT token in the Authorization header
		req.Header.Set("Authorization", user.token)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Perform the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to perform request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != user.status {
			t.Errorf("unexpected status code: %d", resp.StatusCode)
		}

	}

	assert.Nil(t, stop())
}
