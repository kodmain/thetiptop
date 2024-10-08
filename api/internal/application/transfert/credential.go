package transfert

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

type Credential struct {
	ID       *string `json:"id" xml:"id" form:"id"`
	Email    *string `json:"email" xml:"email" form:"email"`
	Password *string `json:"password" xml:"password" form:"password"`
}

func (c *Credential) Check(validator data.Validator) errors.ErrorInterface {
	return validator.Check(data.Object{
		"email":    c.Email,
		"password": c.Password,
	})
}

func NewCredential(obj data.Object, mandatory data.Validator) (*Credential, error) {
	if obj == nil {
		return nil, errors.ErrNoData
	}

	c := &Credential{}

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
