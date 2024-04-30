package validator_test

import (
	"strings"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "Abc123!@#",
			wantErr:  false,
		},
		{
			name:     "Password is too short",
			password: "Abc12",
			wantErr:  true,
		},
		{
			name:     "Password is too long",
			password: "Abc123!@#" + strings.Repeat("a", 57),
			wantErr:  true,
		},
		{
			name:     "Password does not include lowercase letters",
			password: "ABC123!@#",
			wantErr:  true,
		},
		{
			name:     "Password does not include uppercase letters",
			password: "abc123!@#",
			wantErr:  true,
		},
		{
			name:     "Password does not include numbers",
			password: "Abcdef!@#",
			wantErr:  true,
		},
		{
			name:     "Password does not include special characters",
			password: "Abc123456",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Password(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Valid email",
			email:   "hello@kodmain.com",
			wantErr: false,
		},
		{
			name:    "Invalid email",
			email:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Email(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
