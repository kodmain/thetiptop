package transfert_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		email    *string
		password *string
		wantErr  bool
	}{
		{
			name:     "Valid client",
			email:    aws.String("hello@kodmain.com"),
			password: aws.String("Abc123!@#"),
			wantErr:  false,
		},
		{
			name:     "Invalid client",
			email:    aws.String("invalid"),
			password: aws.String(""),
			wantErr:  true,
		},
	}

	client, err := transfert.NewClient(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, client)

	client, err = transfert.NewClient(data.Object{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := data.Object{
				"email":    tt.email,
				"password": tt.password,
			}

			client, err := transfert.NewClient(obj, data.Validator{
				"email":    {validator.Email},
				"password": {validator.Password},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.email, client.Email)
				assert.Equal(t, tt.password, client.Password)
			}
		})
	}
}
