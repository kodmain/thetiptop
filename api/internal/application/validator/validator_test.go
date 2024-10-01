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
		password *string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: aws.String("Abc123!@#"),
			wantErr:  false,
		},
		{
			name:     "Password is too short",
			password: aws.String("Abc12"),
			wantErr:  true,
		},
		{
			name:     "Password is too long",
			password: aws.String("Abc123!@#" + strings.Repeat("a", 57)),
			wantErr:  true,
		},
		{
			name:     "Password does not include lowercase letters",
			password: aws.String("ABC123!@#"),
			wantErr:  true,
		},
		{
			name:     "Password does not include uppercase letters",
			password: aws.String("abc123!@#"),
			wantErr:  true,
		},
		{
			name:     "Password does not include numbers",
			password: aws.String("Abcdef!@#"),
			wantErr:  true,
		},
		{
			name:     "Password does not include special characters",
			password: aws.String("Abc123456"),
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Password(tt.password, "password")
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
		email   *string
		wantErr bool
	}{
		{
			name:    "Valid email",
			email:   aws.String("hello@kodmain.com"),
			wantErr: false,
		},
		{
			name:    "Invalid email",
			email:   aws.String("invalid"),
			wantErr: true,
		},
		{
			name:    "Empty email",
			email:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Email(tt.email, "email")
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
		value   *string
		wantErr bool
	}{
		{
			name:    "Valid Luhn",
			value:   aws.String("378282246310005"),
			wantErr: false,
		},
		{
			name:    "Invalid Luhn",
			value:   aws.String("123456789012345"),
			wantErr: true,
		},
		{
			name:    "Empty Luhn",
			value:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Luhn(tt.value, "token")
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
		uuid    *string
		wantErr bool
	}{
		{
			name:    "Valid UUID",
			uuid:    aws.String("00000000-0000-0000-0000-000000000000"),
			wantErr: false,
		},
		{
			name:    "Invalid UUID",
			uuid:    aws.String("00000000-0000-0000-0000-00000000000"),
			wantErr: true,
		},
		{
			name:    "Empty UUID",
			uuid:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ID(tt.uuid, "id")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsBool(t *testing.T) {
	tests := []struct {
		name    string
		boolean *bool
		wantErr bool
	}{
		{
			name:    "Valid true boolean",
			boolean: aws.Bool(true),
			wantErr: false,
		},
		{
			name:    "Valid false boolean",
			boolean: aws.Bool(false),
			wantErr: false,
		},
		{
			name:    "Empty boolean",
			boolean: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.IsBool(tt.boolean, "boolean")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsTrue(t *testing.T) {
	tests := []struct {
		name    string
		boolean *bool
		wantErr bool
	}{
		{
			name:    "Valid true boolean",
			boolean: aws.Bool(true),
			wantErr: false,
		},
		{
			name:    "Invalid false boolean",
			boolean: aws.Bool(false),
			wantErr: true,
		},
		{
			name:    "Empty boolean",
			boolean: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.IsTrue(tt.boolean, "boolean")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsFalse(t *testing.T) {
	tests := []struct {
		name    string
		boolean *bool
		wantErr bool
	}{

		{
			name:    "Invalid true boolean",
			boolean: aws.Bool(true),
			wantErr: true,
		},
		{
			name:    "Valid false boolean",
			boolean: aws.Bool(false),
			wantErr: false,
		},
		{
			name:    "Empty boolean",
			boolean: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.IsFalse(tt.boolean, "boolean")
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
			err := validator.Required(tt.value, "value")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
