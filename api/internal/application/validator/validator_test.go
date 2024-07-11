package validator_test

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
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
			err := validator.Password(&tt.password)
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
			err := validator.Email(&tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLuhn(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "Valid Luhn",
			value:   "378282246310005",
			wantErr: false,
		},
		{
			name:    "Invalid Luhn",
			value:   "123456789012345",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Luhn(&tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestID(t *testing.T) {
	tests := []struct {
		name    string
		uuid    string
		wantErr bool
	}{
		{
			name:    "Valid UUID",
			uuid:    "00000000-0000-0000-0000-000000000000",
			wantErr: false,
		},
		{
			name:    "Invalid UUID",
			uuid:    "00000000-0000-0000-0000-00000000000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ID(&tt.uuid)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name    string
		value   *string
		wantErr bool
	}{
		{
			name:    "Value is not nil",
			value:   aws.String("not nil"),
			wantErr: false,
		},
		{
			name:    "Value is nil",
			value:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Required(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
