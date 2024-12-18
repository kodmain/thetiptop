package transfert

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

type Store struct {
	ID       *string `json:"id" xml:"id" form:"id"`
	Label    *string `json:"label" xml:"label" form:"label"`
	IsOnline *bool   `json:"is_online" xml:"is_online" form:"is_online"`
}

func (c *Store) Check(validator data.Validator) errors.ErrorInterface {
	return validator.Check(data.Object{
		"id":        c.ID,
		"label":     c.Label,
		"is_online": c.IsOnline,
	})
}

func NewStore(obj data.Object, mandatory data.Validator) (*Store, error) {
	if obj == nil {
		return nil, errors.ErrNoData
	}

	c := &Store{}

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
