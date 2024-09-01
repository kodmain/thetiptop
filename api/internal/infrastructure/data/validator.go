package data

type Control func(any, string) error
type Validator map[string][]Control

func (d Validator) Check(obj Object) error {
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
