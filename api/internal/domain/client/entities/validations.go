package entities

type Validations []*Validation

func (vs Validations) Has(t ValidationType) *Validation {
	for _, v := range vs {
		if v.Type == t {
			return v
		}
	}

	return nil
}
