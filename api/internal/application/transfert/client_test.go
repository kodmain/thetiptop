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
		name       string
		newsletter *bool
		cgu        *bool
		wantErr    bool
	}{
		{
			name:       "Valid client",
			newsletter: aws.Bool(false),
			cgu:        aws.Bool(true),
			wantErr:    false,
		},
		{
			name:       "Invalid client",
			newsletter: aws.Bool(true),
			cgu:        aws.Bool(false),
			wantErr:    true,
		},
	}

	// Test with nil object and nil validator
	client, err := transfert.NewClient(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, client)

	// Test with empty object and nil validator
	client, err = transfert.NewClient(data.Object{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Iterate through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := data.Object{
				"newsletter": tt.newsletter,
				"cgu":        tt.cgu,
			}

			client, err := transfert.NewClient(obj, data.Validator{
				"newsletter": {validator.IsFalse},
				"cgu":        {validator.IsTrue},
			})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)

				// Validate client with the same validators
				err := client.Check(data.Validator{
					"newsletter": {validator.IsFalse},
					"cgu":        {validator.IsTrue},
				})
				assert.NoError(t, err)
			}
		})
	}
}
