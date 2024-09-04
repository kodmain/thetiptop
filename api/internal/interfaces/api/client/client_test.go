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
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
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
//
// Parameters:
// - values: map[string][]string Les valeurs du formulaire à convertir
//
// Returns:
// - url.Values: Les valeurs du formulaire encodées
func createFormValues(values ...map[string][]any) url.Values {
	form := url.Values{}
	if len(values) > 0 {
		for key, valueSlice := range values[0] {
			if len(valueSlice) > 0 {
				switch v := valueSlice[0].(type) {
				case string:
					form.Set(key, v)
				case int:
					form.Set(key, fmt.Sprintf("%d", v))
				case bool:
					form.Set(key, fmt.Sprintf("%t", v))
				default:
					// Gérer les types non supportés
					form.Set(key, fmt.Sprintf("%v", v))
				}
			}
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

// convertFormToJSON convertit url.Values en une map[string]string pour l'encodage JSON
//
// Parameters:
// - form: url.Values Les valeurs du formulaire à convertir
//
// Returns:
// - map[string]string: Les valeurs converties en map pour l'encodage JSON
func convertFormToJSON(form map[string][]any) map[string]any {
	jsonMap := make(map[string]any)
	for key, values := range form {
		if len(values) > 0 {
			jsonMap[key] = values[0] // Garder le type original (string, int, bool)
		}
	}
	return jsonMap
}

// createRequest crée une requête HTTP en fonction des paramètres fournis
//
// Parameters:
// - method: string La méthode HTTP (GET, POST, etc.)
// - uri: string L'URL de la requête
// - token: string Le jeton d'autorisation (si nécessaire)
// - form: url.Values Les valeurs du formulaire encodées
// - encoding: EncodingType Le type d'encodage des données
//
// Returns:
// - *http.Request: La requête HTTP créée
// - error: L'erreur rencontrée (le cas échéant)
func createRequest(method, uri, token string, form map[string][]any, encoding EncodingType) (*http.Request, error) {
	var req *http.Request
	var err error

	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if encoding == JSONEncoded {
			jsonMap := convertFormToJSON(form)
			jsonData, err := json.Marshal(jsonMap)
			if err != nil {
				return nil, err
			}
			req, err = http.NewRequest(method, uri, bytes.NewBuffer(jsonData))
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/json")
		} else {
			formValues := createFormValues(form)
			req, err = http.NewRequest(method, uri, strings.NewReader(formValues.Encode()))
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	case http.MethodGet, http.MethodDelete:
		formValues := createFormValues(form)
		uri = fmt.Sprintf("%s?%s", uri, formValues.Encode())
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
// - encoding: EncodingType Le type d'encodage des données
// - values: map[string][]string Les valeurs du formulaire (facultatif)
//
// Returns:
// - []byte: Le contenu de la réponse
// - int: Le code de statut HTTP
// - error: L'erreur rencontrée (le cas échéant)
func request(method, uri string, token string, encoding EncodingType, values ...map[string][]any) ([]byte, int, error) {
	form := map[string][]any{}
	if len(values) > 0 {
		form = values[0]
	}
	req, err := createRequest(method, uri, token, form, encoding)
	if err != nil {
		return nil, 0, err
	}

	// Effectuer la requête
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
	encodingTypes := []EncodingType{FormURLEncoded, JSONEncoded}
	assert.Nil(t, start(8888, 8444))

	const (
		DOMAIN = "http://localhost:8888"

		USER_REGISTER            = DOMAIN + "/user/register"
		USER_REGISTER_VALIDATION = DOMAIN + "/user/register/validation"
		USER_AUTH                = DOMAIN + "/user/auth"
		USER_AUTH_RENEW          = DOMAIN + "/user/auth/renew"
		USER_PASSWORD_UPDATE     = DOMAIN + "/user/password/update"

		RECOVER_PASSWORD   = DOMAIN + "/recover/password"
		RECOVER_VALIDATION = DOMAIN + "/recover/validation"
	)

	for _, encoding := range encodingTypes {

		var encodingName string = "FormURLEncoded"
		if encoding == JSONEncoded {
			encodingName = "JSONEncoded"
		}

		request("DELETE", "http://0.0.0.0:1080/email/all", "", encoding)
		time.Sleep(1 * time.Second)

		users := []struct {
			email    string
			password string
			statusSU int
			statusSI int
		}{
			// mail, pass, status-signup, status-signin
			{fmt.Sprintf("%v", encoding) + GOOD_EMAIL, GOOD_PASS, http.StatusCreated, http.StatusOK},
			{fmt.Sprintf("%v", encoding) + GOOD_EMAIL, GOOD_PASS + "hello", http.StatusConflict, http.StatusBadRequest},
			{fmt.Sprintf("%v", encoding) + WRONG_EMAIL, WRONG_PASS, http.StatusBadRequest, http.StatusBadRequest},
		}

		t.Run("SignUp/"+encodingName, func(t *testing.T) {
			for _, user := range users {
				values := map[string][]any{
					"oki": {user.email},
				}

				RegisteredClient, status, err := request("POST", USER_REGISTER, "", encoding, values)

				assert.Nil(t, err)
				assert.Equal(t, http.StatusBadRequest, status)

				values = map[string][]any{
					"email":      {user.email},
					"password":   {user.password},
					"newsletter": {true},
					"cgu":        {true},
				}

				RegisteredClient, status, err = request("POST", USER_REGISTER, "", encoding, values)
				assert.Nil(t, err)
				assert.Equal(t, user.statusSU, status)

				if status == http.StatusCreated {
					t.Run("Validation/"+encodingName, func(t *testing.T) {
						var client entities.Client
						err = json.Unmarshal(RegisteredClient, &client)
						assert.NoError(t, err)
						assert.NotNil(t, client)
						time.Sleep(3 * time.Second)
						email, err := getMailFor(user.email)
						assert.Nil(t, err)
						assert.Equal(t, user.email, email.To[0].Address)
					})

					t.Run("Validation/recover/"+encodingName, func(t *testing.T) {
						_, status, err := request("POST", RECOVER_VALIDATION, "", encoding, map[string][]any{
							"email": {user.email},
							"type":  {entities.MailValidation.String()},
						})
						assert.Nil(t, err)
						assert.Equal(t, http.StatusNoContent, status)
						time.Sleep(3 * time.Second)
						email, err := getMailFor(user.email)
						assert.Nil(t, err)
						assert.Equal(t, user.email, email.To[0].Address)
						token := extractToken(email.HTML)
						_, status, err = request("PUT", USER_REGISTER_VALIDATION, "", encoding, map[string][]any{
							"token": {token},
							"email": {user.email},
						})
						assert.Nil(t, err)
						assert.Equal(t, http.StatusOK, status)
					})
				}

				JWT, status, err := request("POST", USER_AUTH, "", encoding, values)
				assert.Nil(t, err)
				assert.Equal(t, user.statusSI, status)

				if status == http.StatusOK {
					t.Run("Renew/"+encodingName, func(t *testing.T) {
						var tokenData fiber.Map
						err = json.Unmarshal(JWT, &tokenData)
						assert.Nil(t, err)
						refresh_token_sting := tokenData["refresh_token"].(string)
						users := []struct {
							token  string
							status int
						}{
							{"Bearer " + refresh_token_sting, http.StatusOK},
							{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDkxMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJ0eXBlIjoxLCJ0eiI6IkxvY2FsIn0.5Lae56HNcQ1OHcP_FhTfcOOtHpaZVgRFy6vzzBugN7Y", http.StatusUnauthorized}, // Replace with actual expired JWT token
							{"Bearer malformed.jwt.token.here", http.StatusUnauthorized},
							{"", http.StatusBadRequest},
						}

						for _, user := range users {
							_, status, err := request("GET", USER_AUTH_RENEW, user.token, encoding)
							assert.Nil(t, err)
							assert.Equal(t, user.status, status)
						}
					})

					t.Run("Password/"+encodingName, func(t *testing.T) {
						_, status, err := request("POST", RECOVER_PASSWORD, "", encoding, map[string][]any{
							"email": {user.email + "wrong"},
						})

						assert.NoError(t, err)
						assert.Equal(t, http.StatusNotFound, status)

						_, status, err = request("POST", RECOVER_PASSWORD, "", encoding, map[string][]any{
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

						_, status, err = request("PUT", USER_PASSWORD_UPDATE, "", encoding, map[string][]any{
							"token":    {token},
							"email":    {user.email},
							"password": {GOOD_PASS_UPDATED},
						})

						assert.Nil(t, err)
						assert.Equal(t, http.StatusOK, status)

						_, status, err = request("POST", USER_AUTH, "", encoding, values)
						assert.NoError(t, err)
						assert.Equal(t, http.StatusBadRequest, status)

						values["password"] = []any{GOOD_PASS_UPDATED}
						_, status, err = request("POST", USER_AUTH, "", encoding, values)
						assert.NoError(t, err)
						assert.Equal(t, user.statusSI, status)
					})
				}

			}
		})
	}
	assert.Nil(t, stop())
}
