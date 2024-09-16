package transfert

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

type Employee struct {
	ID           *string `json:"id" xml:"id" form:"id"`
	Newsletter   *bool   `json:"newsletter" xml:"newsletter" form:"newsletter"`
	CGU          *bool   `json:"cgu" xml:"cgu" form:"cgu"`
	CredentialID *string `json:"credential_id" xml:"credential_id" form:"credential_id"`
}

func (e *Employee) Check(validator data.Validator) error {
	return validator.Check(data.Object{
		"id":            e.ID,
		"newsletter":    e.Newsletter,
		"cgu":           e.CGU,
		"credential_id": e.CredentialID,
	})
}

func NewEmployee(obj data.Object, mandatory data.Validator) (*Employee, error) {
	if obj == nil {
		return nil, fmt.Errorf(errors.ErrNoData)
	}

	c := &Employee{}

	if mandatory == nil {
		if err := obj.Hydrate(c); err != nil {
			return nil, err
		}

		return c, nil
	}

	if err := mandatory.Check(obj); err != nil {
		return nil, err
	}

	if err := obj.Hydrate(c); err != nil {
		return nil, err
	}

	return c, nil
}
