package services_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/security"
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

// TestUserAuth tests the UserAuth function
// This function checks the user authentication process
//
// Parameters:
// - t: *testing.T test framework
//
// Returns:
// - none
func TestUserAuth(t *testing.T) {
	// Load config for tests
	err := config.Load(aws.String("../../../../config.test.yml"))
	assert.NoError(t, err)

	email := "test@example.com"
	password := "ValidPassword123!"
	emailSyntaxFail := "invalid-email"
	passwordSyntaxFail := "short"

	t.Run("invalid syntax password", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		statusCode, response := services.UserAuth(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &passwordSyntaxFail,
		})

		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("invalid syntax email", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		statusCode, response := services.UserAuth(mockClient, &transfert.Credential{
			Email:    &emailSyntaxFail,
			Password: &password,
		})

		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Simulate client not found
		mockClient.On("UserAuth", mock.Anything).
			Return(nil, "", errors_domain_user.ErrClientNotFound)

		statusCode, response := services.UserAuth(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &password,
		})

		assert.Equal(t, fiber.StatusNotFound, statusCode)

		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errObj)
			assert.Equal(t, errors_domain_user.ErrClientNotFound, errObj)
		}

		assert.NotNil(t, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("valid email, password", func(t *testing.T) {
		t.Parallel()

		id, errGen := uuid.NewRandom()
		assert.NoError(t, errGen)

		ids := id.String()
		mockClient := new(DomainUserService)
		// Simulate a successful auth: returns an ID and a role
		mockClient.On("UserAuth", mock.Anything).
			Return(&ids, security.ROLE_CONNECTED, nil)

		statusCode, response := services.UserAuth(mockClient, &transfert.Credential{
			Email:    &email,
			Password: &password,
		})

		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)

		// Because we return fiber.Map{"access_token", ... "refresh_token", ...}
		respMap, ok := response.(fiber.Map)
		assert.True(t, ok, "response should be of type fiber.Map")
		if ok {
			assert.Contains(t, respMap, "access_token")
			assert.Contains(t, respMap, "refresh_token")
		}

		mockClient.AssertExpectations(t)
	})
}

// TestUserAuthRenew tests the UserAuthRenew function
// This function checks the refresh token renewal logic
//
// Parameters:
// - t: *testing.T test framework
//
// Returns:
// - none
func TestUserAuthRenew(t *testing.T) {
	err := config.Load(aws.String("../../../../config.test.yml"))
	assert.NoError(t, err)

	t.Run("invalid token - nil", func(t *testing.T) {
		t.Parallel()

		// Null token
		statusCode, response := services.UserAuthRenew(nil)

		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errObj)
		}
	})

	t.Run("invalid token type", func(t *testing.T) {
		t.Parallel()

		// Token with wrong type
		invalidToken := &jwt.Token{
			Type: jwt.ACCESS, // WRONG type
		}

		statusCode, response := services.UserAuthRenew(invalidToken)
		assert.Equal(t, fiber.StatusUnauthorized, statusCode)

		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errObj)
		}
	})

	t.Run("token expired", func(t *testing.T) {
		t.Parallel()

		// Refresh token that has expired
		expiredToken := &jwt.Token{
			Type: jwt.REFRESH,
			Exp:  time.Now().Add(-1 * time.Hour).Unix(),
		}

		statusCode, response := services.UserAuthRenew(expiredToken)
		assert.Equal(t, fiber.StatusUnauthorized, statusCode)

		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errObj)
		}
	})

	t.Run("successful token renewal", func(t *testing.T) {
		t.Parallel()

		// Valid refresh token
		validToken := &jwt.Token{
			Type: jwt.REFRESH,
			ID:   "valid-client-id",
			Exp:  time.Now().Add(1 * time.Hour).Unix(),
		}

		statusCode, response := services.UserAuthRenew(validToken)
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)

		respMap, ok := response.(fiber.Map)
		assert.True(t, ok, "response should be of type fiber.Map")
		if ok {
			assert.NotNil(t, respMap["access_token"])
			assert.NotNil(t, respMap["refresh_token"])
		}
	})
}

// TestCredentialUpdate tests the CredentialUpdate function
// This function checks the logic for updating credentials after a validation
//
// Parameters:
// - t: *testing.T test framework
//
// Returns:
// - none
func TestCredentialUpdate(t *testing.T) {
	config.Load(aws.String("../../../config.test.yml"))

	luhn := token.Generate(6)

	t.Run("invalid token syntax", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		validationDTO := &transfert.Validation{
			Token: aws.String("invalidToken"),
		}
		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		errMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errMap, "token")
		}
		mockClient.AssertExpectations(t)
	})

	t.Run("invalid email format", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		validationDTO := &transfert.Validation{
			Token: luhn.PointerString(),
		}
		credentialDTO := &transfert.Credential{
			Email:    aws.String("invalid-email"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		errMap, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errMap, "email")
		}
		mockClient.AssertExpectations(t)
	})

	t.Run("password validation error", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// The PasswordValidation returns an error
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).
			Return(nil, errors_domain_user.ErrValidationNotFound)

		validationDTO := &transfert.Validation{
			Token: luhn.PointerString(),
		}
		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationNotFound, response)

		mockClient.AssertExpectations(t)
	})

	t.Run("validation already validated", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// The PasswordValidation returns ErrValidationAlreadyValidated
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).
			Return(nil, errors_domain_user.ErrValidationAlreadyValidated)

		validationDTO := &transfert.Validation{
			Token: luhn.PointerString(),
		}
		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationAlreadyValidated, response)

		mockClient.AssertExpectations(t)
	})

	t.Run("validation expired", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// The PasswordValidation returns ErrValidationExpired
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).
			Return(nil, errors_domain_user.ErrValidationExpired)

		validationDTO := &transfert.Validation{
			Token: luhn.PointerString(),
		}
		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusGone, statusCode)
		assert.Equal(t, errors_domain_user.ErrValidationExpired, response)

		mockClient.AssertExpectations(t)
	})

	t.Run("successful password update", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// Simulate successful password validation and update
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).
			Return(&entities.Validation{}, nil)
		mockClient.On("PasswordUpdate", mock.Anything).Return(nil)

		validationDTO := &transfert.Validation{
			Token: luhn.PointerString(),
		}
		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)

		mockClient.AssertExpectations(t)
	})

	t.Run("error during password update", func(t *testing.T) {
		t.Parallel()

		mockClient := new(DomainUserService)
		// PasswordValidation ok, but PasswordUpdate fails
		mockClient.On("PasswordValidation", mock.Anything, mock.Anything).
			Return(&entities.Validation{}, nil)
		mockClient.On("PasswordUpdate", mock.Anything).
			Return(errors.ErrInternalServer)

		validationDTO := &transfert.Validation{
			Token: luhn.PointerString(),
		}
		credentialDTO := &transfert.Credential{
			Email:    aws.String("test@example.com"),
			Password: aws.String("ValidP@ssw0rd"),
		}

		statusCode, response := services.CredentialUpdate(mockClient, validationDTO, credentialDTO)
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)

		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errObj)
			assert.Equal(t, errors.ErrInternalServer, errObj)
		}

		mockClient.AssertExpectations(t)
	})
}
