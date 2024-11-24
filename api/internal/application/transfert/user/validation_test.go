package transfert_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
)

func TestNewValidation(t *testing.T) {

	luhn := token.Generate(6)
	fakeLuhn := token.NewLuhn("123456")

	tests := []struct {
		name     string
		clientID *string
		luhn     *string
		wantErr  bool
	}{
		{
			name:     "Valid validation",
			clientID: aws.String("00000000-0000-0000-0000-000000000000"),
			luhn:     luhn.PointerString(),
			wantErr:  false,
		},
		{
			name:     "Invalid validation",
			clientID: aws.String("invalid"),
			luhn:     fakeLuhn.PointerString(),
			wantErr:  true,
		},
	}

	validation, err := transfert.NewValidation(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, validation)

	validation, err = transfert.NewValidation(data.Object{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, validation)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := data.Object{
				"client_id": tt.clientID,
				"token":     tt.luhn,
			}

			validation, err := transfert.NewValidation(obj, data.Validator{
				"client_id": {validator.ID},
				"token":     {validator.Luhn},
			})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, validation)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, validation)

				err = validation.Check(data.Validator{
					"client_id": {validator.ID},
					"token":     {validator.Luhn},
				})

				assert.NoError(t, err)
			}
		})
	}

}
