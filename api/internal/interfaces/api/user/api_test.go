package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application/hook"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	userRepository "github.com/kodmain/thetiptop/api/internal/domain/user/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/server"
	"github.com/kodmain/thetiptop/api/internal/interfaces"
)

const (
	GOOD_EMAIL        = "user1@example.com"
	GOOD_PASS         = "ValidP@ssw0rd1"
	GOOD_PASS_UPDATED = "ValidP@ssw0rd1-update"

	WRONG_EMAIL = "user2@example.com"
	WRONG_PASS  = "secret"

	email    = "user-thetiptop@yopmail.com"
	password = "Aa1@azetyuiop"
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

var callBack hook.HandlerSync = func(tags ...string) {
	if len(tags) > 0 && tags[0] == "default" {
		user := userRepository.NewUserRepository(database.Get(config.GetString("services.game.database", config.DEFAULT)))
		cred, _ := user.CreateCredential(&transfert.Credential{
			Email:    aws.String(email),
			Password: aws.String(password),
		})

		user.CreateClient(&transfert.Client{
			CredentialID: &cred.ID,
		})
	}
}

func start(http, https int) error {
	env.DEFAULT_PORT_HTTP = http
	env.DEFAULT_PORT_HTTPS = https
	env.PORT_HTTP = &env.DEFAULT_PORT_HTTP
	env.PORT_HTTPS = &env.DEFAULT_PORT_HTTPS
	env.ForceTest()
	hook.Register(hook.EventOnDBInit, callBack)
	config.Load(aws.String("../../../../config.test.yml"))
	logger.Info("starting application")
	srv = server.Create()
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
		if len(form) > 0 {
			formValues := createFormValues(form)
			uri = fmt.Sprintf("%s?%s", uri, formValues.Encode())
		}

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

// getMailFor Try to retrieve an email for the given address X times before returning an error
// This function will attempt to find and delete an email associated with the given address.
// It will retry the operation X times before failing.
//
// Parameters:
// - emailAddr: string The email address to search for
// - retries: int The number of times to retry before failing
//
// Returns:
// - email: *Email The found email
// - error: error Error if the operation fails after X attempts
func getMailFor(emailAddr string, retries int) (*Email, error) {
	apiUrl := "http://localhost:1080/email"

	for i := 0; i < retries; i++ {
		content, statusCode, err := request("GET", apiUrl, "", FormURLEncoded)

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

		time.Sleep(100 * time.Millisecond)
	}

	return nil, fmt.Errorf("aucun email trouvé pour l'adresse %s après %d tentatives", emailAddr, retries)
}

func extractToken(html string) string {
	re := regexp.MustCompile(`<h1>([a-zA-Z0-9]+)</h1>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

const (
	DOMAIN = "http://localhost:8888"

	// Client
	CLIENT          = DOMAIN + "/client"
	CLIENT_REGISTER = CLIENT + "/register"
	CLIENT_WITH_ID  = CLIENT + "/%s"

	// Employee
	EMPLOYEE          = DOMAIN + "/employee"
	EMPLOYEE_REGISTER = EMPLOYEE + "/register"
	EMPLOYEE_WITH_ID  = EMPLOYEE + "/%s"

	// User
	USER                     = DOMAIN + "/user"
	USER_AUTH                = USER + "/auth"
	USER_AUTH_RENEW          = USER + "/auth/renew"
	USER_PASSWORD            = USER + "/password"
	USER_REGISTER_VALIDATION = USER + "/register/validation"
	USER_VALIDATION_RENEW    = USER + "/validation/renew"
)
