package services_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/services"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	email              = "test@example.com"
	emailSyntaxFail    = "testexample.com"
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

type DomainClientService struct {
	mock.Mock
}

func (dcs DomainClientService) SignUp(obj *transfert.Client) (*entities.Client, error) {
	args := dcs.Called(obj)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (dcs DomainClientService) SignIn(obj *transfert.Client) (*entities.Client, error) {
	args := dcs.Called(obj)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Client), args.Error(1)
}

func (dcs DomainClientService) ValidationMail(obj *transfert.Validation) (*entities.Validation, error) {
	args := dcs.Called(obj)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func TestClient(t *testing.T) {
	assert.Nil(t, setup())

	t.Run("invalid password", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignUp", mock.Anything).Return(&entities.Client{}, nil)

		statusCode, response := services.SignUp(mockClient, email, passwordSyntaxFail)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("valid password", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignUp", mock.Anything).Return(&entities.Client{}, nil)

		statusCode, response := services.SignUp(mockClient, email, password)
		assert.Equal(t, fiber.StatusCreated, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("client already exists", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignUp", mock.Anything).Return(nil, fmt.Errorf(errors.ErrClientAlreadyExists))

		statusCode, response := services.SignUp(mockClient, email, password)
		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.NotNil(t, response)
	})
}

func TestSignIn(t *testing.T) {
	config.Load(aws.String("../../../config.test.yml"))
	t.Run("invalid syntax password", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignIn", mock.Anything).Return(&entities.Client{}, nil)

		statusCode, response := services.SignIn(mockClient, email, passwordSyntaxFail)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("invalid syntax email", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignIn", mock.Anything).Return(&entities.Client{}, nil)

		statusCode, response := services.SignIn(mockClient, emailSyntaxFail, password)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("not found", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignIn", mock.Anything).Return(nil, fmt.Errorf(errors.ErrClientNotFound))

		statusCode, response := services.SignIn(mockClient, email, password)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignIn", mock.Anything).Return(nil, fmt.Errorf("fail to log in"))

		statusCode, response := services.SignIn(mockClient, email, passwordFail)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("valid email,password", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignIn", mock.Anything).Return(&entities.Client{
			ID: "7c79400f-006a-475e-97b6-5dbde3707601",
		}, nil)

		statusCode, response := services.SignIn(mockClient, email, password)
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)

		tokenJWT, ok := response["jwt"].(string)
		assert.True(t, ok)

		access, err := jwt.TokenToClaims(tokenJWT)
		assert.Nil(t, err)
		assert.NotNil(t, access)

		statusCode, response = services.SignRenew(access)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
		assert.Equal(t, "Invalid token type", response["error"])

		assert.False(t, access.HasExpired()) // Le jeton ne doit pas être expiré
		assert.NotNil(t, access.Refresh)     // Le jeton doit avoir un jeton de rafraîchissement

		refresh, err := jwt.TokenToClaims(*access.Refresh)
		if err != nil {
			t.Error(err)
		}

		assert.False(t, refresh.HasExpired())
		assert.Nil(t, refresh.Refresh)
		statusCode, response = services.SignRenew(refresh)
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)

		expired, err := jwt.TokenToClaims(ExpiredRefreshToken)
		assert.Error(t, err)

		statusCode, response = services.SignRenew(expired)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})
}
