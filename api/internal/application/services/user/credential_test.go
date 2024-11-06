package services_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/config"
	services "github.com/kodmain/thetiptop/api/internal/application/services/user"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	errors_domain_user "github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserAuth(t *testing.T) {
	err := config.Load(aws.String("../../../../config.test.yml"))
	assert.NoError(t, err)

	t.Run("invalid syntax password", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Mocker correctement la méthode UserAuth
		mockClient.On("UserAuth", mock.Anything).Return(&entities.Client{}, nil)

		statusCode, response := services.UserAuth(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &passwordSyntaxFail,
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("invalid syntax email", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Mocker correctement la méthode UserAuth
		mockClient.On("UserAuth", mock.Anything).Return(&entities.Client{}, nil)

		statusCode, response := services.UserAuth(mockClient, &transfert.Credential{
			Email:    &emailSyntaxFail,
			Password: &password,
		})
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("not found", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Simuler le cas où le client n'est pas trouvé
		mockClient.On("UserAuth", mock.Anything).Return(nil, "", errors_domain_user.ErrClientNotFound)

		statusCode, response := services.UserAuth(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &password,
		})
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Error(t, response.(*errors.Error))
		assert.NotNil(t, response)
	})

	t.Run("valid email, password", func(t *testing.T) {
		id, err := uuid.NewRandom()
		ids := id.String()
		assert.NoError(t, err)
		mockClient := new(DomainUserService)
		// Simuler un cas réussi avec une Credential valide et un ClientID valide
		mockClient.On("UserAuth", mock.Anything).Return(&ids, "", nil)

		statusCode, response := services.UserAuth(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &password,
		})
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)
	})
}

func TestUserAuthRenew(t *testing.T) {
	err := config.Load(aws.String("../../../../config.test.yml"))
	assert.NoError(t, err)

	t.Run("invalid token - nil", func(t *testing.T) {
		// Cas où le jeton est nil
		statusCode, response := services.UserAuthRenew(nil)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("invalid token type", func(t *testing.T) {
		// Cas où le type de jeton est invalide
		invalidToken := &jwt.Token{
			Type: jwt.ACCESS, // Mauvais type de jeton
		}

		statusCode, response := services.UserAuthRenew(invalidToken)
		assert.Equal(t, fiber.StatusUnauthorized, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("token expired", func(t *testing.T) {
		// Cas où le jeton de rafraîchissement a expiré
		expiredToken := &jwt.Token{
			Type: jwt.REFRESH,
			Exp:  time.Now().Add(-1 * time.Hour).Unix(), // Jeton expiré
		}

		statusCode, response := services.UserAuthRenew(expiredToken)
		assert.Equal(t, fiber.StatusUnauthorized, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("successful token renewal", func(t *testing.T) {
		// Cas de renouvellement réussi avec un jeton valide
		validToken := &jwt.Token{
			Type: jwt.REFRESH,
			ID:   "valid-client-id",
			Exp:  time.Now().Add(1 * time.Hour).Unix(), // Jeton valide
		}

		statusCode, response := services.UserAuthRenew(validToken)
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)

		// Vérifier que les jetons sont présents dans la réponse
		authResponse, ok := response.(fiber.Map)
		assert.True(t, ok)
		assert.NotNil(t, authResponse["access_token"])
		assert.NotNil(t, authResponse["refresh_token"])
	})
}

func TestCredentialUpdate(t *testing.T) {
	config.Load(aws.String("../../../config.test.yml"))

	luhn := token.Generate(6)

	t.Run("invalid token syntax", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Cas où le token est invalide
		validationDTO := &transfert.Validation{
			Token: aws.String("invalidToken"),
		}

		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("invalid email format", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Cas où l'email est invalide
		validationDTO := &transfert.Validation{
			Token: luhn.PointerString(),
		}

		credentialDTO := &transfert.Credential{
			Email:    aws.String("invalid-email"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("password validation error", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Simuler une erreur lors de la validation du mot de passe
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).Return(nil, errors_domain_user.ErrValidationNotFound)

		validationDTO := &transfert.Validation{
			Token: luhn.Pointer().PointerString(),
		}

		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationNotFound, response)
	})

	t.Run("validation already validated", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Simuler le cas où la validation a déjà été effectuée
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).Return(nil, errors_domain_user.ErrValidationAlreadyValidated)

		validationDTO := &transfert.Validation{
			Token: luhn.Pointer().PointerString(),
		}

		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationAlreadyValidated, response)
	})

	t.Run("validation expired", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Simuler le cas où la validation a expiré
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).Return(nil, errors_domain_user.ErrValidationExpired)

		validationDTO := &transfert.Validation{
			Token: luhn.Pointer().PointerString(),
		}

		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusGone, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationExpired, response)
	})

	t.Run("successful password update", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Simuler une mise à jour réussie du mot de passe
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).Return(&entities.Validation{}, nil)
		mockClient.On("PasswordUpdate", mock.Anything).Return(nil)

		validationDTO := &transfert.Validation{
			Token: luhn.Pointer().PointerString(),
		}

		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)
	})

	t.Run("error during password update", func(t *testing.T) {
		mockClient := new(DomainUserService)
		// Simuler une erreur lors de la mise à jour du mot de passe
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).Return(&entities.Validation{}, nil)
		mockClient.On("PasswordUpdate", mock.Anything).Return(errors.ErrInternalServer)

		validationDTO := &transfert.Validation{
			Token: luhn.Pointer().PointerString(),
		}

		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Error(t, response.(*errors.Error))
	})
}
