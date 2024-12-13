package transfert

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

type Caisse struct {
	ID      *string `json:"id" xml:"id" form:"id"`
	Label   *string `json:"label" xml:"label" form:"label"`
	StoreID *string `json:"store_id" xml:"store_id" form:"store_id"`
}

func (c *Caisse) Check(validator data.Validator) errors.ErrorInterface {
	return validator.Check(data.Object{
		"id":       c.ID,
		"label":    c.Label,
		"store_id": c.StoreID,
	})
}

func NewCaisse(obj data.Object, mandatory data.Validator) (*Caisse, error) {
	if obj == nil {
		return nil, errors.ErrNoData
	}

	c := &Caisse{}

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
