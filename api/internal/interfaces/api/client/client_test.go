package client_test

import (
	"bytes"
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
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/server"
	"github.com/kodmain/thetiptop/api/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

const (
	GOOD_EMAIL        = "user1@example.com"
	GOOD_PASS         = "ValidP@ssw0rd1"
	GOOD_PASS_UPDATED = "ValidP@ssw0rd1-update"

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

// EncodingType représente le type d'encodage des données
type EncodingType int

const (
	FormURLEncoded EncodingType = iota
	JSONEncoded
)

func createRequest(method, uri, token string, form url.Values, encoding EncodingType) (*http.Request, error) {
	var req *http.Request
	var err error

	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if encoding == JSONEncoded {
			jsonData, err := json.Marshal(form)
			if err != nil {
				return nil, err
			}
			req, err = http.NewRequest(method, uri, bytes.NewBuffer(jsonData))
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, err = http.NewRequest(method, uri, strings.NewReader(form.Encode()))
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	case http.MethodGet, http.MethodDelete:
		uri = fmt.Sprintf("%s?%s", uri, form.Encode())
		req, err = http.NewRequest(method, uri, nil)
		if err != nil {
			return nil, err
		}
	default:
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

func request(method, uri string, token string, encoding EncodingType, values ...map[string][]string) ([]byte, int, error) {
	form := createFormValues(values...)
	req, err := createRequest(method, uri, token, form, encoding)
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

// deleteMail Delete email by ID
// This function deletes an email from the server by its ID.
//
// Parameters:
// - emailID: string ID of the email to be deleted
//
// Returns:
// - err: error Error if any occurred during the deletion process
func deleteMail(emailID string) error {
	apiUrl := fmt.Sprintf("http://localhost:1080/email/%s", emailID)
	_, statusCode, err := request("DELETE", apiUrl, "", FormURLEncoded)

	if err != nil {
		return fmt.Errorf("erreur lors de la requête HTTP: %v", err)
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("statut de réponse non OK: %d", statusCode)
	}

	return nil
}

func getMailFor(emailAddr string) (*Email, error) {
	apiUrl := "http://localhost:1080/email"
	var content []byte
	var statusCode int
	var err error

	content, statusCode, err = request("GET", apiUrl, "", FormURLEncoded)

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
				deleteErr := deleteMail(email.ID)
				if deleteErr != nil {
					return nil, fmt.Errorf("erreur lors de la suppression de l'email: %v", deleteErr)
				}
				return email, nil
			}
		}
	}

	return nil, fmt.Errorf("aucun email trouvé pour l'adresse %s", emailAddr)
}

func generateRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func extractURLs(html string) []string {
	re := regexp.MustCompile(`https?://[^\s"']+`)
	matches := re.FindAllString(html, -1)
	return matches
}

func extractToken(html string) string {
	re := regexp.MustCompile(`<h1>([a-zA-Z0-9]+)</h1>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func TestClient(t *testing.T) {
	assert.Nil(t, start(8888, 8444))

	request("DELETE", "http://0.0.0.0:1080/email/all", "", FormURLEncoded)
	time.Sleep(1 * time.Second)

	users := []struct {
		email    string
		password string
		statusSU int
		statusSI int
	}{
		// mail, pass, status-signup, status-signin
		{GOOD_EMAIL, GOOD_PASS, http.StatusCreated, http.StatusOK},
		{GOOD_EMAIL, GOOD_PASS + "hello", http.StatusConflict, http.StatusBadRequest},
		{WRONG_EMAIL, WRONG_PASS, http.StatusBadRequest, http.StatusBadRequest},
	}

	t.Run("SignUp", func(t *testing.T) {
		for _, user := range users {
			values := map[string][]string{
				"email":    {user.email},
				"password": {user.password},
			}

			RegisteredClient, status, err := request("POST", "http://localhost:8888/sign/up", "", FormURLEncoded, values)
			assert.Nil(t, err)
			assert.Equal(t, user.statusSU, status)

			if status == http.StatusCreated {
				t.Run("Validation", func(t *testing.T) {
					var clientWrapper map[string]entities.Client
					err = json.Unmarshal(RegisteredClient, &clientWrapper)
					assert.NoError(t, err)
					client := clientWrapper["client"]
					assert.NotNil(t, client)
					time.Sleep(1 * time.Second)
					email, err := getMailFor(user.email)
					assert.Nil(t, err)
					assert.Equal(t, user.email, email.To[0].Address)
					token := extractToken(email.HTML)
					logger.Info(token)
					_, status, err = request("PUT", "http://localhost:8888/sign/validation", "", FormURLEncoded, map[string][]string{
						"token": {token},
						"email": {user.email},
					})
					assert.Nil(t, err)
					assert.Equal(t, http.StatusOK, status)
				})
			}

			JWT, status, err := request("POST", "http://localhost:8888/sign/in", "", FormURLEncoded, values)
			assert.Nil(t, err)
			assert.Equal(t, user.statusSI, status)

			if status == http.StatusOK {
				t.Run("Renew", func(t *testing.T) {
					var tokenData TokenStructure
					err = json.Unmarshal(JWT, &tokenData)
					assert.Nil(t, err)

					access, err := serializer.TokenToClaims(tokenData.JWT)
					assert.Nil(t, err)

					users := []struct {
						token  string
						status int
					}{
						{"Bearer " + *access.Refresh, http.StatusOK},
						{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDkxMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJ0eXBlIjoxLCJ0eiI6IkxvY2FsIn0.5Lae56HNcQ1OHcP_FhTfcOOtHpaZVgRFy6vzzBugN7Y", http.StatusUnauthorized}, // Replace with actual expired JWT token
						{"Bearer malformed.jwt.token.here", http.StatusUnauthorized},
						{"", http.StatusBadRequest},
					}

					for _, user := range users {
						_, status, err := request("GET", "http://localhost:8888/sign/renew", user.token, FormURLEncoded)
						assert.Nil(t, err)
						assert.Equal(t, user.status, status)
					}
				})

				t.Run("Password", func(t *testing.T) {
					_, status, err := request("POST", "http://localhost:8888/password/recover", "", FormURLEncoded, map[string][]string{
						"email": {user.email + "wrong"},
					})

					assert.Error(t, err)
					assert.Equal(t, http.StatusNotFound, status)

					_, status, err = request("POST", "http://localhost:8888/password/recover", "", FormURLEncoded, map[string][]string{
						"email": {user.email},
					})

					assert.Nil(t, err)
					assert.Equal(t, http.StatusNoContent, status)
					time.Sleep(1 * time.Second)
					email, err := getMailFor(user.email)
					assert.Nil(t, err)
					assert.Equal(t, user.email, email.To[0].Address)

					token := extractToken(email.HTML)
					assert.NotEmpty(t, token)

					_, status, err = request("PUT", "http://localhost:8888/password/update", "", FormURLEncoded, map[string][]string{
						"token":    {token},
						"email":    {user.email},
						"password": {GOOD_PASS_UPDATED},
					})

					assert.Nil(t, err)
					assert.Equal(t, http.StatusOK, status)

					_, status, err = request("POST", "http://localhost:8888/sign/in", "", FormURLEncoded, values)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusBadRequest, status)

					values["password"] = []string{GOOD_PASS_UPDATED}
					_, status, err = request("POST", "http://localhost:8888/sign/in", "", FormURLEncoded, values)
					assert.NoError(t, err)
					assert.Equal(t, user.statusSI, status)
				})
			}

		}
	})

	assert.Nil(t, stop())
}
