package transfert_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/stretchr/testify/assert"
)

func TestNewTicket(t *testing.T) {
	tests := []struct {
		name      string
		inputData data.Object
		wantErr   bool
	}{
		{
			name: "Valid ticket",
			inputData: data.Object{
				"id":        aws.String("123"),
				"prize":     aws.String("Gold"),
				"client_id": aws.String("456"),
				"token":     aws.String("abc123"),
			},
			wantErr: false,
		},
		{
			name: "Invalid ticket - missing ID",
			inputData: data.Object{
				"prize":     aws.String("Gold"),
				"client_id": aws.String("456"),
				"token":     aws.String("abc123"),
			},
			wantErr: true,
		},
		{
			name: "Invalid ticket - missing Token",
			inputData: data.Object{
				"id":        aws.String("123"),
				"prize":     aws.String("Gold"),
				"client_id": aws.String("456"),
			},
			wantErr: true,
		},
	}

	// Test for nil object and nil validator
	t.Run("Nil object and validator", func(t *testing.T) {
		ticket, err := transfert.NewTicket(nil, nil)
		assert.Error(t, err, "expected error when object is nil")
		assert.Nil(t, ticket, "expected ticket to be nil")
	})

	// Test for empty object with nil validator
	t.Run("Empty object and nil validator", func(t *testing.T) {
		ticket, err := transfert.NewTicket(data.Object{}, nil)
		assert.NoError(t, err, "expected no error for empty object")
		assert.NotNil(t, ticket, "expected ticket to be created")
	})

	// Iterate through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticket, err := transfert.NewTicket(tt.inputData, data.Validator{
				"id":        {validator.Required},
				"prize":     {validator.Required},
				"client_id": {validator.Required},
				"token":     {validator.Required},
			})

			if tt.wantErr {
				assert.Error(t, err, "expected error for invalid input")
				assert.Nil(t, ticket, "expected ticket to be nil for invalid input")
			} else {
				assert.NoError(t, err, "expected no error for valid input")
				assert.NotNil(t, ticket, "expected valid ticket to be created")

				// Validate the ticket object with the same validators
				err := ticket.Check(data.Validator{
					"id":        {validator.Required},
					"prize":     {validator.Required},
					"client_id": {validator.Required},
					"token":     {validator.Required},
				})
				assert.NoError(t, err, "expected ticket validation to pass")
			}
		})
	}
}
