package transfert_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid email",
			email:    "hello@kodmain.com",
			password: "Abc123!@#",
			wantErr:  false,
		},
		{
			name:     "Invalid email",
			email:    "invalid",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := data.Object{
				"email":    tt.email,
				"password": tt.password,
			}

			client, err := transfert.NewClient(obj)
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
