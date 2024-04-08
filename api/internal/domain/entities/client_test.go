package entities_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/dto"
	"github.com/kodmain/thetiptop/api/internal/architecture/security"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	email := "test@example.com"
	password := "password"

	dto := &dto.Client{
		Email:    email,
		Password: password,
	}

	client, err := entities.CreateClient(dto)
	assert.Nil(t, err)

	if client.ID == uuid.Nil {
		t.Error("Expected non-zero UUID, got zero")
	}

	if client.Email != email {
		t.Errorf("Expected email %s, got %s", email, client.Email)
	}

	hash, err := security.Hash(email+":"+password, security.BCRYPT)
	assert.Nil(t, err)
	if security.CompareHash(client.Password, email+":"+password, security.BCRYPT) != nil {
		t.Errorf("Expected password %s, got %s", hash, client.Password)
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
