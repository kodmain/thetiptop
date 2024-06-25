package entities_test

import (
	"testing"
	"time"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	email := "user-thetiptop@yopmail.com"
	password := "Aa1@azetyuiop"

	obj := &transfert.Client{
		Email:    email,
		Password: password,
	}

	client := entities.CreateClient(obj)

	if client.Email != email {
		t.Errorf("Expected email %s, got %s", email, client.Email)
	}

	hash, err := hash.Hash(email+":"+password, hash.BCRYPT)
	assert.Nil(t, err)
	if client.CompareHash(email + ":" + password) {
		t.Errorf("Expected password %s", hash)
	}

	now := time.Now()
	if client.CreatedAt.After(now) {
		t.Errorf("Expected CreatedAt to be before or equal to current time, got %s", client.CreatedAt)
	}

	if client.UpdatedAt.After(now) {
		t.Errorf("Expected UpdatedAt to be before or equal to current time, got %s", client.UpdatedAt)
	}

	if client.DeletedAt.Valid {
		t.Error("Expected DeletedAt to be nil, got non-nil")
	}
}
func TestBefore(t *testing.T) {
	email := "user-thetiptop@yopmail.com"
	password := "Aa1@azetyuiop"

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
