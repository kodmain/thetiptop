package transfert

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

type Client struct {
	Email      *string `json:"email" xml:"email" form:"email"`
	Password   *string `json:"password" xml:"password" form:"password"`
	Newsletter *string `json:"newsletter" xml:"newsletter" form:"newsletter"`
	CGU        *string `json:"cgu" xml:"cgu" form:"cgu"`
}

func (c *Client) Check(validator data.Validator) error {
	return validator.Check(data.Object{
		"email":      c.Email,
		"password":   c.Password,
		"newsletter": c.Newsletter,
		"cgu":        c.CGU,
	})
}

func NewClient(obj data.Object, mandatory data.Validator) (*Client, error) {
	if obj == nil {
		return nil, fmt.Errorf(errors.ErrNoData)
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
