package transfert_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"

	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

func TestNewCaisse(t *testing.T) {
	tests := []struct {
		name     string
		label    *string
		store_id *string
		wantErr  bool
	}{
		{
			name:     "Valid caisse",
			store_id: aws.String("387f3fb0-88a0-4e2f-bc82-529719e5ed21"),
			label:    aws.String("Caisse 1"),
			wantErr:  false,
		},
		{
			name:     "Invalid caisse",
			store_id: aws.String("true"),
			label:    aws.String("Caisse 2"),
			wantErr:  true,
		},
	}

	// Test with nil object and nil validator
	caisse, err := transfert.NewCaisse(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, caisse)

	// Test with empty object and nil validator
	caisse, err = transfert.NewCaisse(data.Object{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, caisse)

	// Iterate through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := data.Object{
				"store_id": tt.store_id,
				"label":    tt.label,
			}

			caisse, err := transfert.NewCaisse(obj, data.Validator{
				"store_id": {validator.Required, validator.ID},
				"label":    {validator.Required},
			})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, caisse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, caisse)

				// Validate caisse with the same validators
				err := caisse.Check(data.Validator{
					"store_id": {validator.Required, validator.ID},
					"label":    {validator.Required},
				})
				assert.NoError(t, err)
			}
		})
	}
}
