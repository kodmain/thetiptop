package transfert

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

type Validation struct {
	ID         *string `json:"id" xml:"id" form:"id"`
	Token      *string `json:"token" xml:"token" form:"token"`
	ClientID   *string `json:"client_id" xml:"client_id" form:"client_id"`
	EmployeeID *string `json:"employee_id" xml:"employee_id" form:"employee_id"`
	Type       *string `json:"type" xml:"type" form:"type"`
}

func (v *Validation) Check(validator data.Validator) error {
	return validator.Check(data.Object{
		"id":          v.ID,
		"token":       v.Token,
		"client_id":   v.ClientID,
		"employee_id": v.EmployeeID,
		"type":        v.Type,
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
