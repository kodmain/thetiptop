package entities_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"github.com/stretchr/testify/assert"
)

func TestClient_HasSuccessValidation(t *testing.T) {
	client := &entities.Client{
		Validations: entities.Validations{
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: true,
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: false,
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: true,
			},
		},
	}

	validationType := entities.PasswordRecover
	result := client.HasSuccessValidation(entities.PasswordRecover)

	assert.NotNil(t, result)
	assert.Equal(t, validationType, result.Type)
	assert.True(t, result.Validated)

	result = client.HasSuccessValidation(entities.MailValidation)
	assert.Nil(t, result)

}

func TestClient_HasNotExpiredValidation(t *testing.T) {
	client := &entities.Client{
		Validations: entities.Validations{
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: true,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: false,
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: false,
				ExpiresAt: time.Now().Add(time.Hour),
			},
		},
	}

	validationType := entities.PasswordRecover
	result := client.HasNotExpiredValidation(entities.PasswordRecover)

	assert.NotNil(t, result)
	assert.Equal(t, validationType, result.Type)
	assert.False(t, result.Validated)
	assert.False(t, result.HasExpired())

	result = client.HasNotExpiredValidation(entities.MailValidation)
	assert.Nil(t, result)
}

func TestNewClient(t *testing.T) {
	email := aws.String("user-thetiptop@yopmail.com")
	password := aws.String("Aa1@azetyuiop")

	obj := &transfert.Client{
		Email:    email,
		Password: password,
	}

	client := entities.CreateClient(obj)

	if *client.Email != *email {
		t.Errorf("Expected email %s, got %s", *email, *client.Email)
	}

	hash, err := hash.Hash(aws.String(*email+":"+*password), hash.BCRYPT)
	assert.Nil(t, err)
	if client.CompareHash(*email + ":" + *password) {
		t.Errorf("Expected password %s", *hash)
	}

	now := time.Now()
	if client.CreatedAt.After(now) {
		t.Errorf("Expected CreatedAt to be before or equal to current time, got %s", client.CreatedAt)
	}

	if client.UpdatedAt.After(now) {
		t.Errorf("Expected UpdatedAt to be before or equal to current time, got %s", client.UpdatedAt)
	}
}
func TestClientBefore(t *testing.T) {
	email := aws.String("user-thetiptop@yopmail.com")
	password := aws.String("Aa1@azetyuiop")

	client := &entities.Client{
		Email:    email,
		Password: password,
	}

	err := client.BeforeCreate(nil)
	assert.Nil(t, err)

	assert.NotEmpty(t, client.ID)
	assert.NotEmpty(t, client.Password)

	old := client.UpdatedAt
	time.Sleep(100 * time.Millisecond)

	err = client.BeforeUpdate(nil)
	assert.Nil(t, err)
	assert.True(t, client.UpdatedAt.After(old))
}
