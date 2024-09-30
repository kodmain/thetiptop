package data

import "github.com/kodmain/thetiptop/api/internal/infrastructure/errors"

type Control func(any, string) errors.ErrorInterface
type Validator map[string][]Control

func (d Validator) Check(obj Object) errors.ErrorInterface {
	for key, controls := range d {
		value := obj.Get(key)
		for _, control := range controls {
			if err := control(value, key); err != nil {
				return err
			}
		}
	}

	return nil
}
