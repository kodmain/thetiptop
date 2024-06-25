package entities

import "encoding/json"

// ValidationType defines the type of validation
type ValidationType int

const (
	Mail ValidationType = iota
	Phone
)

var validationTypeToString = map[ValidationType]string{
	Mail:  "email",
	Phone: "phone",
}

var stringToValidationType = map[string]ValidationType{
	"email": Mail,
	"phone": Phone,
}

// MarshalJSON marshals the enum as a string
func (v ValidationType) MarshalJSON() ([]byte, error) {
	return json.Marshal(validationTypeToString[v])
}

// UnmarshalJSON unmarshals a string to the enum type
func (v *ValidationType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*v = stringToValidationType[str]
	return nil
}
