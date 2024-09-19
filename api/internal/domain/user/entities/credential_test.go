package entities_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"github.com/stretchr/testify/assert"
)

func TestNewCredential(t *testing.T) {
	email := aws.String("user-thetiptop@yopmail.com")
	password := aws.String("Aa1@azetyuiop")

	obj := &transfert.Credential{
		Email:    email,
		Password: password,
	}

	credential := entities.CreateCredential(obj)

	if *credential.Email != *email {
		t.Errorf("Expected email %s, got %s", *email, *credential.Email)
	}

	hash, err := hash.Hash(aws.String(*email+":"+*password), hash.BCRYPT)
	assert.Nil(t, err)
	if credential.CompareHash(*email + ":" + *password) {
		t.Errorf("Expected password %s", *hash)
	}

	now := time.Now()
	if credential.CreatedAt.After(now) {
		t.Errorf("Expected CreatedAt to be before or equal to current time, got %s", credential.CreatedAt)
	}

	if credential.UpdatedAt.After(now) {
		t.Errorf("Expected UpdatedAt to be before or equal to current time, got %s", credential.UpdatedAt)
	}
}

func TestCredentialBefore(t *testing.T) {
	email := aws.String("user-thetiptop@yopmail.com")
	password := aws.String("Aa1@azetyuiop")

	credential := &entities.Credential{
		Email:    email,
		Password: password,
	}

	err := credential.BeforeCreate(nil)
	assert.Nil(t, err)

	assert.NotEmpty(t, credential.ID)
	assert.NotEmpty(t, credential.Password)

	old := credential.UpdatedAt
	time.Sleep(100 * time.Millisecond)

	err = credential.BeforeUpdate(nil)
	assert.Nil(t, err)
	assert.True(t, credential.UpdatedAt.After(old))
}
