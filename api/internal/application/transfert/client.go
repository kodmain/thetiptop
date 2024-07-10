package transfert

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

type Client struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
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
