package transfert_test

import (
	"testing"

	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/stretchr/testify/assert"
)

func TestNewEmployee(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Valid employee",
			wantErr: false,
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
			obj := data.Object{}
			employee, err := transfert.NewEmployee(obj, data.Validator{})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, employee)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, employee)
				err := employee.Check(data.Validator{})
				assert.NoError(t, err)
			}
		})
	}
}
