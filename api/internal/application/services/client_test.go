package services_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/services"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	email              = "test@example.com"
	emailSyntaxFail    = "testexample.com"
	password           = "validP@ssw0rd"
	passwordFail       = "WrongP@ssw0rd"
	passwordSyntaxFail = "secret"

	ExpiredAccessToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDgyMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJyZWZyZXNoIjoiZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFM01UTXhNRGt4TXpFc0ltbGtJam9pTjJNM09UUXdNR1l0TURBMllTMDBOelZsTFRrM1lqWXROV1JpWkdVek56QTNOakF4SWl3aWIyWm1Jam8zTWpBd0xDSjBlWEJsSWpveExDSjBlaUk2SWt4dlkyRnNJbjAuNUxhZTU2SE5jUTFPSGNQX0ZoVGZjT090SHBhWlZnUkZ5NnZ6ekJ1Z043WSIsInR5cGUiOjAsInR6IjoiTG9jYWwifQ.BxW2wfHiiCr0aTsuWwRVmh0Wd-BX20AoUDTGg_rIDoM"
	ExpiredRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDkxMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJ0eXBlIjoxLCJ0eiI6IkxvY2FsIn0.5Lae56HNcQ1OHcP_FhTfcOOtHpaZVgRFy6vzzBugN7Y"
)

type DomainClientService struct {
	mock.Mock
}

func (dcs DomainClientService) PasswordRecover(obj *transfert.Client) error {
	args := dcs.Called(obj)
	return args.Error(0)
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

func (dcs DomainClientService) SignValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Client) (*entities.Validation, error) {
	args := dcs.Called(dtoValidation, dtoClient)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (dcs DomainClientService) PasswordValidation(dtoValidation *transfert.Validation, dtoClient *transfert.Client) (*entities.Validation, error) {
	args := dcs.Called(dtoValidation, dtoClient)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Validation), args.Error(1)
}

func (dcs DomainClientService) PasswordUpdate(client *transfert.Client) error {
	args := dcs.Called(client)
	return args.Error(0)
}

func TestClient(t *testing.T) {
	config.Load(aws.String("../../../config.test.yml"))

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

	t.Run("client fail to signup", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignUp", mock.Anything).Return(nil, fmt.Errorf("boom"))

		statusCode, response := services.SignUp(mockClient, email, password)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
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
		id, err := uuid.Parse("7c79400f-006a-475e-97b6-5dbde3707601")
		assert.NoError(t, err)
		mockClient := new(DomainClientService)
		mockClient.On("SignIn", mock.Anything).Return(&entities.Client{
			ID: id.String(),
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

	t.Run("client fail to signin", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignIn", mock.Anything).Return(nil, fmt.Errorf("boom"))

		statusCode, response := services.SignIn(mockClient, email, password)
		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})
}

func TestSignValidation(t *testing.T) {
	config.Load(aws.String("../../../config.test.yml"))
	luhn := token.Generate(6)

	t.Run("invalid syntax token", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignValidation", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("invalid token"))

		statusCode, response := services.SignValidation(mockClient, email, "invalidToken")
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
		assert.Equal(t, "invalid digit", response["error"])
	})

	t.Run("valid token, email", func(t *testing.T) {
		mockClient := new(DomainClientService)
		id, err := uuid.NewRandom()
		assert.NoError(t, err)

		mockClient.On("SignValidation", mock.Anything, mock.Anything).Return(&entities.Validation{
			ID: id.String(),
		}, nil)

		statusCode, response := services.SignValidation(mockClient, email, luhn.String())
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)
	})
	t.Run("validation not found", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignValidation", mock.Anything, mock.Anything).Return(nil, fmt.Errorf(errors.ErrValidationNotFound))

		statusCode, response := services.SignValidation(mockClient, email, luhn.String())
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.NotNil(t, response)
		assert.Equal(t, errors.ErrValidationNotFound, response["error"])
	})

	t.Run("validation already validated", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignValidation", mock.Anything, mock.Anything).Return(nil, fmt.Errorf(errors.ErrValidationAlreadyValidated))

		statusCode, response := services.SignValidation(mockClient, email, luhn.String())
		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.NotNil(t, response)
		assert.Equal(t, errors.ErrValidationAlreadyValidated, response["error"])
	})

	t.Run("validation already validated", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("SignValidation", mock.Anything, mock.Anything).Return(nil, fmt.Errorf(errors.ErrValidationExpired))

		statusCode, response := services.SignValidation(mockClient, email, luhn.String())
		assert.Equal(t, fiber.StatusGone, statusCode)
		assert.NotNil(t, response)
		assert.Equal(t, errors.ErrValidationExpired, response["error"])
	})
}

func TestPasswordRecover(t *testing.T) {
	config.Load(aws.String("../../../config.test.yml"))

	t.Run("invalid syntax email", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("PasswordRecover", mock.Anything).Return(fmt.Errorf("invalid email"))

		statusCode, response := services.PasswordRecover(mockClient, emailSyntaxFail)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("valid email", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("PasswordRecover", mock.Anything).Return(nil)

		statusCode, response := services.PasswordRecover(mockClient, email)
		assert.Equal(t, fiber.StatusNoContent, statusCode)
		assert.Nil(t, response)
	})

	t.Run("client fail to recover password", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("PasswordRecover", mock.Anything).Return(fmt.Errorf("boom"))

		statusCode, response := services.PasswordRecover(mockClient, email)
		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("client not found", func(t *testing.T) {
		mockClient := new(DomainClientService)
		mockClient.On("PasswordRecover", mock.Anything).Return(fmt.Errorf(errors.ErrClientNotFound))

		statusCode, response := services.PasswordRecover(mockClient, email)
		assert.Equal(t, http.StatusNotFound, statusCode)
		assert.NotNil(t, response)
	})
}

func TestPasswordUpdate(t *testing.T) {
	config.Load(aws.String("../../../config.test.yml"))

	t.Run("invalid syntax email", func(t *testing.T) {
		config.Load(aws.String("../../../config.test.yml"))

		mockClient := new(DomainClientService)
		mockClient.On("PasswordUpdate", mock.Anything).Return(fmt.Errorf("invalid email"))

		statusCode, response := services.PasswordUpdate(mockClient, emailSyntaxFail, password, "token")
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("invalid syntax password", func(t *testing.T) {
		config.Load(aws.String("../../../config.test.yml"))

		mockClient := new(DomainClientService)
		mockClient.On("PasswordUpdate", mock.Anything).Return(fmt.Errorf("invalid password"))

		statusCode, response := services.PasswordUpdate(mockClient, email, passwordSyntaxFail, "token")
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("invalid token", func(t *testing.T) {
		config.Load(aws.String("../../../config.test.yml"))

		mockClient := new(DomainClientService)
		mockClient.On("PasswordUpdate", mock.Anything).Return(fmt.Errorf("invalid password"))

		statusCode, response := services.PasswordUpdate(mockClient, email, password, "token")
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("validation not found", func(t *testing.T) {
		config.Load(aws.String("../../../config.test.yml"))

		mockClient := new(DomainClientService)
		mockClient.On("PasswordUpdate", mock.Anything).Return(nil)
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).Return(nil, fmt.Errorf(errors.ErrValidationNotFound))

		luhn := token.Generate(6)

		statusCode, response := services.PasswordUpdate(mockClient, email, password, luhn.String())
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("validation already validate", func(t *testing.T) {
		config.Load(aws.String("../../../config.test.yml"))

		mockClient := new(DomainClientService)
		mockClient.On("PasswordUpdate", mock.Anything).Return(nil)
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).Return(nil, fmt.Errorf(errors.ErrValidationAlreadyValidated))

		luhn := token.Generate(6)

		statusCode, response := services.PasswordUpdate(mockClient, email, password, luhn.String())
		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("validation not found", func(t *testing.T) {
		config.Load(aws.String("../../../config.test.yml"))

		mockClient := new(DomainClientService)
		mockClient.On("PasswordUpdate", mock.Anything).Return(nil)
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).Return(nil, fmt.Errorf(errors.ErrValidationExpired))

		luhn := token.Generate(6)

		statusCode, response := services.PasswordUpdate(mockClient, email, password, luhn.String())
		assert.Equal(t, fiber.StatusGone, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("valid email, password, token", func(t *testing.T) {
		config.Load(aws.String("../../../config.test.yml"))

		mockClient := new(DomainClientService)
		mockClient.On("PasswordUpdate", mock.Anything).Return(nil)
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).Return(&entities.Validation{}, nil)

		luhn := token.Generate(6)

		statusCode, response := services.PasswordUpdate(mockClient, email, password, luhn.String())
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("client fail to update password", func(t *testing.T) {
		config.Load(aws.String("../../../config.test.yml"))

		mockClient := new(DomainClientService)
		mockClient.On("PasswordUpdate", mock.Anything).Return(fmt.Errorf("boom"))
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).Return(&entities.Validation{}, nil)

		luhn := token.Generate(6)

		statusCode, response := services.PasswordUpdate(mockClient, email, password, luhn.String())
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.NotNil(t, response)
	})

}
