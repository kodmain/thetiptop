package services_test

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/services"
	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
)

const (
	email              = "test@example.com"
	password           = "validP@ssw0rd"
	passwordFail       = "WrongP@ssw0rd"
	passwordSyntaxFail = "secret"

	ExpiredAccessToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDgyMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJyZWZyZXNoIjoiZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM01UTXhNRGt4TXpFc0ltbGtJam9pTjJNM09UUXdNR1l0TURBMllTMDBOelZsTFRrM1lqWXROV1JpWkdVek56QTNOakF4SWl3aWIyWm1Jam8zTWpBd0xDSjBlWEJsSWpveExDSjBlaUk2SWt4dlkyRnNJbjAuNUxhZTU2SE5jUTFPSGNQX0ZoVGZjT090SHBhWlZnUkZ5NnZ6ekJ1Z043WSIsInR5cGUiOjAsInR6IjoiTG9jYWwifQ.BxW2wfHiiCr0aTsuWwRVmh0Wd-BX20AoUDTGg_rIDoM"
	ExpiredRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDkxMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJ0eXBlIjoxLCJ0eiI6IkxvY2FsIn0.5Lae56HNcQ1OHcP_FhTfcOOtHpaZVgRFy6vzzBugN7Y"
)

func setup() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = config.Load(aws.String(workingDir + "/../../../config.test.yml"))
	if err != nil {
		return err
	}

	return nil
}

func TestClient(t *testing.T) {
	assert.Nil(t, setup())

	// Test de l'inscription avec un mot de passe invalide
	statusCode, response := services.SignUp(email, passwordSyntaxFail)
	assert.Equal(t, fiber.StatusBadRequest, statusCode)
	assert.NotNil(t, response)

	// Test de la première inscription avec un mot de passe valide
	statusCode, response = services.SignUp(email, password)
	assert.Equal(t, fiber.StatusCreated, statusCode)
	assert.Nil(t, response)

	// Test de la tentative de réinscription avec le même email
	statusCode, response = services.SignUp(email, password)
	assert.Equal(t, fiber.StatusConflict, statusCode)
	assert.NotNil(t, response)

	// Test de la connexion avec un mot de passe incorrect
	statusCode, response = services.SignIn(email, passwordSyntaxFail)
	assert.Equal(t, fiber.StatusBadRequest, statusCode)
	assert.NotNil(t, response, "should fail to log in")

	// Test de la connexion avec un mot de passe incorrect
	statusCode, response = services.SignIn(email, passwordFail)
	assert.Equal(t, fiber.StatusBadRequest, statusCode)
	assert.NotNil(t, response, "should fail to log in")

	// Test de connexion avec le même utilisateur
	statusCode, response = services.SignIn(email, password)
	assert.Equal(t, fiber.StatusOK, statusCode)
	assert.NotNil(t, response, "should successfully log in")

	jwt, ok := response["jwt"].(string)
	if !ok {
		t.Error("JWT token is missing")
	}

	access, err := serializer.TokenToClaims(jwt)
	if err != nil {
		t.Error(err)
	}

	statusCode, response = services.SignRenew(access)
	assert.Equal(t, fiber.StatusBadRequest, statusCode)
	assert.NotNil(t, response)
	assert.Equal(t, "Invalid token type", response["error"])

	assert.False(t, access.HasExpired()) // Le jeton ne doit pas être expiré
	assert.NotNil(t, access.Refresh)     // Le jeton doit avoir un jeton de rafraîchissement

	refresh, err := serializer.TokenToClaims(*access.Refresh)
	if err != nil {
		t.Error(err)
	}

	assert.False(t, refresh.HasExpired())
	assert.Nil(t, refresh.Refresh)
	statusCode, response = services.SignRenew(refresh)
	assert.Equal(t, fiber.StatusOK, statusCode)
	assert.NotNil(t, response)

	expired, err := serializer.TokenToClaims(ExpiredRefreshToken)
	assert.Error(t, err)

	statusCode, response = services.SignRenew(expired)
	assert.Equal(t, fiber.StatusBadRequest, statusCode)
	assert.NotNil(t, response)
}
