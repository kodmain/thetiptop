package transfert

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

type Validation struct {
	Token    *string `json:"token" xml:"token" form:"token"`
	ClientID *string `json:"client_id" xml:"client_id" form:"client_id"`
}

func (v *Validation) Check(validator data.Validator) error {
	return validator.Check(data.Object{
		"token":     v.Token,
		"client_id": v.ClientID,
	})
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
