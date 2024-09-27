package transfert

import (
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

type Employee struct {
	ID           *string `json:"id" xml:"id" form:"id"`
	CredentialID *string `json:"credential_id" xml:"credential_id" form:"credential_id"`
}

func (e *Employee) Check(validator data.Validator) errors.ErrorInterface {
	return validator.Check(data.Object{
		"id":            e.ID,
		"credential_id": e.CredentialID,
	})
}

func NewEmployee(obj data.Object, mandatory data.Validator) (*Employee, error) {
	if obj == nil {
		return nil, errors.ErrNoData
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
