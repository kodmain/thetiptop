package transfert

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
)

type Validation struct {
	Token    *token.Luhn `json:"token"`
	ClientID *string     `json:"client_id"`
}

func NewValidation(obj data.Object, mandatory data.Validator) (*Validation, error) {
	if obj == nil {
		return nil, fmt.Errorf(errors.ErrNoData)
	}

	v := &Validation{}

	if mandatory == nil {
		if err := obj.Hydrate(v); err != nil {
			return nil, err
		}

		return v, nil
	}

	if err := mandatory.Check(obj); err != nil {
		return nil, err
	}

	if err := obj.Hydrate(v); err != nil {
		return nil, err
	}

	return v, nil
}
