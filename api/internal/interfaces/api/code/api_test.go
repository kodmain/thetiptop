package code_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application/hook"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/server"
	"github.com/kodmain/thetiptop/api/internal/interfaces"
)

var srv *server.Server

func start(http, https int) error {
	env.ForceTest()
	hook.Reset()
	env.DEFAULT_PORT_HTTP = http
	env.DEFAULT_PORT_HTTPS = https
	env.PORT_HTTP = &env.DEFAULT_PORT_HTTP
	env.PORT_HTTPS = &env.DEFAULT_PORT_HTTPS
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

const (
	DOMAIN = "http://localhost:8888"

	// Error
	CODE       = DOMAIN + "/code"
	CODE_ERROR = CODE + "/error"
)
