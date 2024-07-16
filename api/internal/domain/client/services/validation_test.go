package services_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
)

func TestSignValidation(t *testing.T) {
	email := aws.String("test@example.com")
	luhn := token.Generate(6)
	now := time.Now()
	notExpired := now.Add(time.Hour)
	expired := now.Add(-time.Hour)

	dtoClient := &transfert.Client{
		Email: email,
	}

	dtoValidation := &transfert.Validation{
		Token: luhn.PointerString(),
	}

	client := &entities.Client{
		ID:    "1",
		Email: dtoClient.Email,
	}

	t.Run("SignValidation", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			Validation := &entities.Validation{
				ClientID:  &client.ID,
				Token:     &luhn,
				Type:      entities.MailValidation,
				ExpiresAt: notExpired,
				Validated: false,
			}

			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(client, nil)
			mockRepo.On("ReadValidation", dtoValidation).Return(Validation, nil)
			mockRepo.On("UpdateValidation", Validation).Return(nil)

			result, err := service.SignValidation(dtoValidation, dtoClient)
			assert.NoError(t, err)
			assert.Equal(t, Validation, result)

			mockRepo.AssertExpectations(t)
		})

		t.Run("ClientNotFound", func(t *testing.T) {
			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
			result, err := service.SignValidation(dtoValidation, dtoClient)
			assert.Error(t, err)
			assert.Nil(t, result)

			mockRepo.AssertExpectations(t)
		})

		t.Run("ValidationNotFound", func(t *testing.T) {
			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(client, nil)
			mockRepo.On("ReadValidation", dtoValidation).Return(nil, fmt.Errorf(errors.ErrValidationNotFound))
			result, err := service.SignValidation(dtoValidation, dtoClient)
			assert.Error(t, err)
			assert.Nil(t, result)

			mockRepo.AssertExpectations(t)
		})

		t.Run("ValidationExpired", func(t *testing.T) {
			Validation := &entities.Validation{
				ClientID:  &client.ID,
				Token:     &luhn,
				Type:      entities.MailValidation,
				ExpiresAt: expired,
				Validated: false,
			}

			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(client, nil)
			mockRepo.On("ReadValidation", dtoValidation).Return(Validation, nil)
			result, err := service.SignValidation(dtoValidation, dtoClient)
			assert.Error(t, err)
			assert.Nil(t, result)

			mockRepo.AssertExpectations(t)
		})

		t.Run("AlreadyValidate", func(t *testing.T) {
			Validation := &entities.Validation{
				ClientID:  &client.ID,
				Token:     &luhn,
				Type:      entities.MailValidation,
				ExpiresAt: notExpired,
				Validated: true,
			}

			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(client, nil)
			mockRepo.On("ReadValidation", dtoValidation).Return(Validation, nil)

			result, err := service.SignValidation(dtoValidation, dtoClient)
			assert.Error(t, err)
			assert.Nil(t, result)

			mockRepo.AssertExpectations(t)
		})

		t.Run("FailUpdate", func(t *testing.T) {
			Validation := &entities.Validation{
				ClientID:  &client.ID,
				Token:     &luhn,
				Type:      entities.MailValidation,
				ExpiresAt: notExpired,
				Validated: false,
			}

			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(client, nil)
			mockRepo.On("ReadValidation", dtoValidation).Return(Validation, nil)
			mockRepo.On("UpdateValidation", Validation).Return(fmt.Errorf("error"))

			result, err := service.SignValidation(dtoValidation, dtoClient)
			assert.Error(t, err)
			assert.Nil(t, result)

			mockRepo.AssertExpectations(t)
		})

	})

	t.Run("PasswordValidation", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			Validation := &entities.Validation{
				ClientID:  &client.ID,
				Token:     &luhn,
				Type:      entities.PasswordRecover,
				ExpiresAt: notExpired,
				Validated: false,
			}

			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(client, nil)
			mockRepo.On("ReadValidation", dtoValidation).Return(Validation, nil)
			mockRepo.On("UpdateValidation", Validation).Return(nil)

			result, err := service.PasswordValidation(dtoValidation, dtoClient)
			assert.NoError(t, err)
			assert.Equal(t, Validation, result)

			mockRepo.AssertExpectations(t)
		})

		t.Run("ClientNotFound", func(t *testing.T) {
			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(nil, fmt.Errorf(errors.ErrClientNotFound))
			result, err := service.PasswordValidation(dtoValidation, dtoClient)
			assert.Error(t, err)
			assert.Nil(t, result)

			mockRepo.AssertExpectations(t)
		})

		t.Run("ValidationNotFound", func(t *testing.T) {
			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(client, nil)
			mockRepo.On("ReadValidation", dtoValidation).Return(nil, fmt.Errorf(errors.ErrValidationNotFound))
			result, err := service.PasswordValidation(dtoValidation, dtoClient)
			assert.Error(t, err)
			assert.Nil(t, result)

			mockRepo.AssertExpectations(t)
		})

		t.Run("ValidationExpired", func(t *testing.T) {
			Validation := &entities.Validation{
				ClientID:  &client.ID,
				Token:     &luhn,
				Type:      entities.MailValidation,
				ExpiresAt: expired,
				Validated: false,
			}

			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(client, nil)
			mockRepo.On("ReadValidation", dtoValidation).Return(Validation, nil)
			result, err := service.PasswordValidation(dtoValidation, dtoClient)
			assert.Error(t, err)
			assert.Nil(t, result)

			mockRepo.AssertExpectations(t)
		})

		t.Run("AlreadyValidate", func(t *testing.T) {
			Validation := &entities.Validation{
				ClientID:  &client.ID,
				Token:     &luhn,
				Type:      entities.MailValidation,
				ExpiresAt: notExpired,
				Validated: true,
			}

			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(client, nil)
			mockRepo.On("ReadValidation", dtoValidation).Return(Validation, nil)

			result, err := service.PasswordValidation(dtoValidation, dtoClient)
			assert.Error(t, err)
			assert.Nil(t, result)

			mockRepo.AssertExpectations(t)
		})

		t.Run("FailUpdate", func(t *testing.T) {
			Validation := &entities.Validation{
				ClientID:  &client.ID,
				Token:     &luhn,
				Type:      entities.MailValidation,
				ExpiresAt: notExpired,
				Validated: false,
			}

			service, mockRepo, _ := setup()
			mockRepo.On("ReadClient", dtoClient).Return(client, nil)
			mockRepo.On("ReadValidation", dtoValidation).Return(Validation, nil)
			mockRepo.On("UpdateValidation", Validation).Return(fmt.Errorf("error"))

			result, err := service.PasswordValidation(dtoValidation, dtoClient)
			assert.Error(t, err)
			assert.Nil(t, result)

			mockRepo.AssertExpectations(t)
		})

	})
}
