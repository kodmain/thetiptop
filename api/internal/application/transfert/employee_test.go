package transfert_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/stretchr/testify/assert"
)

func TestNewEmployee(t *testing.T) {
	tests := []struct {
		name       string
		newsletter *bool
		cgu        *bool
		wantErr    bool
	}{
		{
			name:       "Valid employee",
			newsletter: aws.Bool(false),
			cgu:        aws.Bool(true),
			wantErr:    false,
		},
		{
			name:       "Invalid employee",
			newsletter: aws.Bool(true),
			cgu:        aws.Bool(false),
			wantErr:    true,
		},
	}

	// Test with nil object and nil validator
	employee, err := transfert.NewEmployee(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, employee)

	// Test with empty object and nil validator
	employee, err = transfert.NewEmployee(data.Object{}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, employee)

	// Iterate through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := data.Object{
				"newsletter": tt.newsletter,
				"cgu":        tt.cgu,
			}

			employee, err := transfert.NewEmployee(obj, data.Validator{
				"newsletter": {validator.IsFalse},
				"cgu":        {validator.IsTrue},
			})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, employee)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, employee)

				// Validate employee with the same validators
				err := employee.Check(data.Validator{
					"newsletter": {validator.IsFalse},
					"cgu":        {validator.IsTrue},
				})
				assert.NoError(t, err)
			}
		})
	}
}
