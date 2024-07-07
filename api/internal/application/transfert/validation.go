package transfert

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
)

type Validation struct {
	Token    *token.Luhn `json:"token"`
	ClientID *string     `json:"client_id"`
}

func NewValidation(obj data.Object, mandatory data.Validator) (*Validation, error) {
	if err := mandatory.Check(obj); err != nil {
		return nil, err
	}

	v := &Validation{}

	if err := obj.Hydrate(v); err != nil {
		return nil, err
	}

	return v, nil
}
