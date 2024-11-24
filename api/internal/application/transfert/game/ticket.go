package transfert

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

type Ticket struct {
	ID           *string `json:"id" xml:"id" form:"id"`
	Prize        *string `json:"prize" xml:"prize" form:"prize"`
	CredentialID *string `json:"credential_id" xml:"credential_id" form:"credential_id"`
	Token        *string `json:"token" xml:"token" form:"token"`
}

func (c *Ticket) Check(validator data.Validator) error {
	return validator.Check(data.Object{
		"id":            c.ID,
		"prize":         c.Prize,
		"credential_id": c.CredentialID,
		"token":         c.Token,
	})
}

func NewTicket(obj data.Object, mandatory data.Validator) (*Ticket, error) {
	if obj == nil {
		return nil, errors.ErrNoData
	}

	c := &Ticket{}

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
