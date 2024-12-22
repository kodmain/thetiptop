package transfert_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"

	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

func TestNewStore(t *testing.T) {
	tests := []struct {
		name      string
		label     *string
		is_online *bool
		wantErr   bool
	}{
		{
			name:      "Valid caisse",
			is_online: aws.Bool(true),
			label:     aws.String("Store 1"),
			wantErr:   false,
		},
		{
			name:      "Invalid caisse",
			is_online: aws.Bool(false),
			label:     nil,
			wantErr:   true,
		},
	}

	// Test with nil object and nil validator
	caisse, err := transfert.NewStore(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, caisse)

	// Test with empty object and nil validator
	caisse, err = transfert.NewStore(data.Object{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, caisse)

	// Iterate through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := data.Object{
				"is_online": tt.is_online,
				"label":     tt.label,
			}

			caisse, err := transfert.NewStore(obj, data.Validator{
				"is_online": {validator.Required},
				"label":     {validator.Required},
			})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, caisse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, caisse)

				// Validate caisse with the same validators
				err := caisse.Check(data.Validator{
					"is_online": {validator.Required},
					"label":     {validator.Required},
				})
				assert.NoError(t, err)
			}
		})
	}
}
