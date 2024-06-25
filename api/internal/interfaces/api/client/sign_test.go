package client_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/server"
	"github.com/kodmain/thetiptop/api/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

const (
	GOOD_EMAIL = "user1@example.com"
	GOOD_PASS  = "ValidP@ssw0rd1"

	WRONG_EMAIL = "user2@example.com"
	WRONG_PASS  = "secret"
)

type Email struct {
	HTML    string `json:"html"`
	Text    string `json:"text"`
	Subject string `json:"subject"`
	From    []struct {
		Address string `json:"address"`
		Name    string `json:"name"`
	} `json:"from"`
	To []struct {
		Address string `json:"address"`
		Name    string `json:"name"`
	} `json:"to"`
	ID       string `json:"id"`
	Time     string `json:"time"`
	Read     bool   `json:"read"`
	Envelope struct {
		From struct {
			Address string `json:"address"`
			Args    struct {
				BODY     string `json:"BODY"`
				SMTPUTF8 bool   `json:"SMTPUTF8"`
			} `json:"args"`
		} `json:"from"`
		To []struct {
			Address string `json:"address"`
			Args    bool   `json:"args"`
		} `json:"to"`
		Host          string `json:"host"`
		RemoteAddress string `json:"remoteAddress"`
	} `json:"envelope"`
	Source        string      `json:"source"`
	Size          int         `json:"size"`
	SizeHuman     string      `json:"sizeHuman"`
	Attachments   interface{} `json:"attachments"`
	CalculatedBcc []struct {
		Address string `json:"address"`
		Name    string `json:"name"`
	} `json:"calculatedBcc"`
}

type TokenStructure struct {
	JWT string `json:"jwt"`
}

var srv *server.Server

func start(http, https int) error {
	env.DEFAULT_PORT_HTTP = http
	env.DEFAULT_PORT_HTTPS = https
	env.PORT_HTTP = &env.DEFAULT_PORT_HTTP
	env.PORT_HTTPS = &env.DEFAULT_PORT_HTTPS
	config.Load(aws.String("../../../../config.test.yml"))
	logger.Info("starting application")
	srv = server.Create()

	logger.Warn(*env.PORT_HTTP)
	srv.Register(interfaces.Endpoints)

	return srv.Start()
}

func stop() error {
	logger.Info("waiting for application to shutdown")
	return srv.Stop()
}

// createFormValues crée les valeurs du formulaire à partir des paramètres fournis
func createFormValues(values ...map[string][]string) url.Values {
	form := url.Values{}
	if len(values) > 0 {
		for key, value := range values[0] {
			form.Set(key, value[0])
		}
	}
	return form
}

// createRequest crée une requête HTTP en fonction de la méthode et des paramètres
func createRequest(method, uri, token string, form url.Values) (*http.Request, error) {
	var req *http.Request
	var err error

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req, err = http.NewRequest(method, uri, strings.NewReader(form.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if method == http.MethodGet || method == http.MethodDelete {
		uri = fmt.Sprintf("%s?%s", uri, form.Encode())
		req, err = http.NewRequest(method, uri, nil)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, uri, nil)
		if err != nil {
			return nil, err
		}
	}

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	return req, nil
}

// request Effectue une requête HTTP avec les paramètres fournis
//
// Parameters:
// - method: string La méthode HTTP (GET, POST, etc.)
// - uri: string L'URL de la requête
// - token: string Le jeton d'autorisation (si nécessaire)
// - values: map[string][]string Les valeurs du formulaire (facultatif)
//
// Returns:
// - []byte: Le contenu de la réponse
// - int: Le code de statut HTTP
// - error: L'erreur rencontrée (le cas échéant)
func request(method, uri string, token string, values ...map[string][]string) ([]byte, int, error) {
	form := createFormValues(values...)
	req, err := createRequest(method, uri, token, form)
	if err != nil {
		return nil, 0, err
	}

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	return content, resp.StatusCode, err
}

func getMailFor(emailAddr string) (*Email, error) {
	apiUrl := "http://localhost:1080/email"
	var content []byte
	var statusCode int
	var err error

	content, statusCode, err = request("GET", apiUrl, "")

	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête HTTP: %v", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("statut de réponse non OK: %d", statusCode)
	}

	var emails []*Email
	if err := json.Unmarshal(content, &emails); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing des emails: %v", err)
	}

	for _, email := range emails {
		for _, to := range email.To {
			if to.Address == emailAddr {
				return email, nil
			}
		}
	}

	return nil, fmt.Errorf("aucun email trouvé pour l'adresse %s", emailAddr)
}

func TestSignUp(t *testing.T) {
	assert.Nil(t, start(8888, 8444))

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
		values := map[string][]string{
			"email":    {user.email},
			"password": {user.password},
		}

		_, status, err := request("POST", "http://localhost:8888/sign/up", "", values)
		assert.Nil(t, err)
		assert.Equal(t, user.status, status)
	}

	assert.Nil(t, stop())
}

func generateRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func extractURLs(html string) []string {
	re := regexp.MustCompile(`https?://[^\s"']+`)
	matches := re.FindAllString(html, -1)
	return matches
}

func TestValidationMail(t *testing.T) {
	assert.Nil(t, start(8889, 8445))

	EMAIL := fmt.Sprintf("%d", generateRandomNumber(1, 1000)) + GOOD_EMAIL

	_, status, err := request("POST", "http://localhost:8889/sign/up", "", map[string][]string{
		"email":    {EMAIL},
		"password": {GOOD_PASS},
	})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)
	time.Sleep(1 * time.Second)
	email, err := getMailFor(EMAIL)
	assert.Nil(t, err)
	assert.Equal(t, EMAIL, email.To[0].Address)

	urls := extractURLs(email.HTML)
	_, status, err = request("GET", urls[0], "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	assert.Nil(t, stop())
}

func TestSignIn(t *testing.T) {
	assert.Nil(t, start(8890, 8446))

	EMAIL := fmt.Sprintf("%d", generateRandomNumber(1, 1000)) + GOOD_EMAIL

	request("POST", "http://localhost:8890/sign/up", "", map[string][]string{
		"email":    {EMAIL},
		"password": {GOOD_PASS},
	})

	time.Sleep(1 * time.Second)
	email, err := getMailFor(EMAIL)
	assert.Nil(t, err)
	assert.Equal(t, EMAIL, email.To[0].Address)

	urls := extractURLs(email.HTML)
	_, status, err := request("GET", urls[0], "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	users := []struct {
		email    string
		password string
		status   int
	}{
		{EMAIL, GOOD_PASS, http.StatusOK},
		{EMAIL, WRONG_PASS, http.StatusBadRequest},
		{GOOD_PASS, WRONG_PASS, http.StatusBadRequest},
	}

	for _, user := range users {
		values := map[string][]string{
			"email":    {user.email},
			"password": {user.password},
		}

		_, status, err := request("POST", "http://localhost:8890/sign/in", "", values)
		assert.Nil(t, err)
		assert.Equal(t, user.status, status)
	}

	assert.Nil(t, stop())
}

func TestSignRenew(t *testing.T) {
	assert.Nil(t, start(8891, 8447))

	EMAIL := fmt.Sprintf("%d", generateRandomNumber(1, 1000)) + GOOD_EMAIL

	request("POST", "http://localhost:8891/sign/up", "", map[string][]string{
		"email":    {EMAIL},
		"password": {GOOD_PASS},
	})

	time.Sleep(1 * time.Second)
	email, err := getMailFor(EMAIL)
	assert.Nil(t, err)
	assert.Equal(t, EMAIL, email.To[0].Address)

	urls := extractURLs(email.HTML)
	_, status, err := request("GET", urls[0], "")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	content, _, _ := request("POST", "http://localhost:8891/sign/in", "", map[string][]string{
		"email":    {EMAIL},
		"password": {GOOD_PASS},
	})

	// Déclaration de la variable qui recevra la valeur désérialisée
	var tokenData TokenStructure

	// Désérialisation du JSON dans la structure définie
	err = json.Unmarshal(content, &tokenData)
	assert.Nil(t, err)

	access, err := serializer.TokenToClaims(tokenData.JWT)
	assert.Nil(t, err)

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
		_, status, err := request("GET", "http://localhost:8891/sign/renew", user.token)
		assert.Nil(t, err)
		assert.Equal(t, user.status, status)
	}

	assert.Nil(t, stop())
}
