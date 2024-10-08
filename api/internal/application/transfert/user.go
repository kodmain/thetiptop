package transfert

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

type User struct {
	ID           *string `json:"id" xml:"id" form:"id"`
	CredentialID *string `json:"credential_id" xml:"credential_id" form:"credential_id"`
}

func (u *User) ToClient() *Client {
	return &Client{
		ID:           u.ID,
		CredentialID: u.CredentialID,
	}
}

func (u *User) ToEmployee() *Employee {
	return &Employee{
		ID:           u.ID,
		CredentialID: u.CredentialID,
	}
}

func (c *User) Check(validator data.Validator) error {
	return validator.Check(data.Object{
		"id":            c.ID,
		"credential_id": c.CredentialID,
	})
}

func NewUser(obj data.Object, mandatory data.Validator) (*User, error) {
	if obj == nil {
		return nil, errors.ErrNoData
	}

	c := &User{}

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
