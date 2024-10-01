package transfert_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	uuid := aws.String("00000000-0000-0000-0000-000000000000")

	tests := []struct {
		name         string
		credentialID *string
		id           *string
		validator    data.Validator
		wantErr      bool
	}{
		{
			name:         "Valid user",
			credentialID: uuid,
			id:           uuid,
			validator: data.Validator{
				"id":            {validator.ID},
				"credential_id": {validator.ID},
			},
			wantErr: false,
		},
		{
			name:         "Invalid user (missing id)",
			credentialID: uuid,
			id:           nil,
			validator: data.Validator{
				"id":            {validator.ID},
				"credential_id": {validator.ID},
			},
			wantErr: true,
		},
		{
			name:         "Invalid user (missing credential_id)",
			credentialID: nil,
			id:           aws.String("user-id-456"),
			validator: data.Validator{
				"id":            {validator.ID},
				"credential_id": {validator.ID},
			},
			wantErr: true,
		},
	}

	// Test with nil object and nil validator
	user, err := transfert.NewUser(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, user)

	// Test with empty object and nil validator
	user, err = transfert.NewUser(data.Object{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Iterate through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := data.Object{
				"id":            tt.id,
				"credential_id": tt.credentialID,
			}

			user, err := transfert.NewUser(obj, tt.validator)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)

				// Validate user with the same validators
				err := user.Check(tt.validator)
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserToClient(t *testing.T) {
	// Initialisation de la structure User
	credentialID := aws.String("cred-id-123")
	id := aws.String("user-id-456")
	user := &transfert.User{
		ID:           id,
		CredentialID: credentialID,
	}

	// Conversion vers Client
	client := user.ToClient()

	// Vérification des résultats
	assert.NotNil(t, client)
	assert.Equal(t, user.ID, client.ID)
	assert.Equal(t, user.CredentialID, client.CredentialID)
}

func TestUserToEmployee(t *testing.T) {
	// Initialisation de la structure User
	credentialID := aws.String("cred-id-123")
	id := aws.String("user-id-456")
	user := &transfert.User{
		ID:           id,
		CredentialID: credentialID,
	}

	// Conversion vers Employee
	employee := user.ToEmployee()

	// Vérification des résultats
	assert.NotNil(t, employee)
	assert.Equal(t, user.ID, employee.ID)
	assert.Equal(t, user.CredentialID, employee.CredentialID)
}
