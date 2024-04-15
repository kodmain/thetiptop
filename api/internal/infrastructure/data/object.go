package data

type Object map[string]string

func (d Object) Get(key string) string {
	if value, ok := d[key]; ok {
		return value
	}

	return ""
}
