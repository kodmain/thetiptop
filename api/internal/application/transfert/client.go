package transfert

import (
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

type Client struct {
	ID           *string `json:"id" xml:"id" form:"id"`
	Newsletter   *bool   `json:"newsletter" xml:"newsletter" form:"newsletter"`
	CGU          *bool   `json:"cgu" xml:"cgu" form:"cgu"`
	CredentialID *string `json:"credential_id" xml:"credential_id" form:"credential_id"`
}

func (c *Client) Check(validator data.Validator) errors.ErrorInterface {
	return validator.Check(data.Object{
		"id":            c.ID,
		"newsletter":    c.Newsletter,
		"cgu":           c.CGU,
		"credential_id": c.CredentialID,
	})
}

func NewClient(obj data.Object, mandatory data.Validator) (*Client, error) {
	if obj == nil {
		return nil, errors.ErrNoData
	}

	c := &Client{}

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
