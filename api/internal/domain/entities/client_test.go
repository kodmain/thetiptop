package entities_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
)

func TestNewClient(t *testing.T) {
	email := "test@example.com"
	password := "password"

	client := entities.NewClient(email, password)

	if client.ID == uuid.Nil {
		t.Error("Expected non-zero UUID, got zero")
	}

	if client.Email != email {
		t.Errorf("Expected email %s, got %s", email, client.Email)
	}

	if client.Password != password {
		t.Errorf("Expected password %s, got %s", password, client.Password)
	}

	now := time.Now()
	if client.CreatedAt.After(now) {
		t.Errorf("Expected CreatedAt to be before or equal to current time, got %s", client.CreatedAt)
	}

	if client.UpdatedAt.After(now) {
		t.Errorf("Expected UpdatedAt to be before or equal to current time, got %s", client.UpdatedAt)
	}

	if client.DeletedAt != nil {
		t.Error("Expected DeletedAt to be nil, got non-nil")
	}
}
