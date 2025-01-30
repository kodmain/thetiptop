package data

import "github.com/kodmain/thetiptop/api/internal/infrastructure/errors"

type Control func(any, string) errors.ErrorInterface
type Validator map[string][]Control

func (d Validator) Check(obj Object) errors.ErrorInterface {
	var errList errors.Errors = make(errors.Errors, 0)

	for key, controls := range d {
		value := obj.Get(key)
		for _, control := range controls {
			if err := control(value, key); err != nil {
				errList.Add(key, err)
			}
		}
	}

	return errList.ToErrorInterface()
}
