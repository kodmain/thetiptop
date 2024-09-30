package data_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

func TestValidator_Check(t *testing.T) {

	ns1 := newString("value1")
	ns2 := newString("value2")
	// setup control functions
	successControl := func(value any, name string) errors.ErrorInterface {
		stringValue, ok := value.(*string)
		if !ok {
			return errors.ErrBadRequest
		}

		if value == nil || *stringValue == "" {
			return errors.ErrBadRequest
		}
		return nil
	}

	// create a Validator instance
	validator := data.Validator{
		"key1": {successControl},
		"key2": {successControl},
	}

	tests := []struct {
		name    string
		obj     data.Object
		wantErr bool
	}{
		{
			name: "all controls pass",
			obj: data.Object{
				"key1": ns1,
				"key2": ns2,
			},
			wantErr: false,
		},
		{
			name: "control fails",
			obj: data.Object{
				"key1": ns1,
				"key2": newString(""),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Check(tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v, for %v", err, tt.wantErr, tt.name)
			}
		})
	}
}

// newString est une fonction utilitaire pour retourner un pointeur vers une nouvelle chaîne
func newString(s string) *string {
	return &s
}
