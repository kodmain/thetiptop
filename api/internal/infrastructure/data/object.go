package data

import "encoding/json"

type Object map[string]any

func (d Object) Get(key string) any {
	if value, ok := d[key]; ok {
		return value
	}

	return nil
}

func (d Object) Hydrate(target any) error {
	bytes, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &target)
}
