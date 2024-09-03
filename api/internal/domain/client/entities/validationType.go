package entities

import (
	"encoding/json"
	"fmt"
)

// ValidationType defines the type of validation
type ValidationType int

const (
	MailValidation ValidationType = iota
	PhoneValidation
	PasswordRecover
)

var validationTypeToString = map[ValidationType]string{
	MailValidation:  "mail",
	PhoneValidation: "phone",
	PasswordRecover: "password",
}

var stringToValidationType = map[string]ValidationType{
	"mail":     MailValidation,
	"phone":    PhoneValidation,
	"password": PasswordRecover,
}

func newValidationType(v *string) (ValidationType, error) {
	if v == nil {
		return 0, fmt.Errorf("ValidationType is required")
	}

	return stringToValidationType[*v], nil
}

func (v ValidationType) String() string {
	return validationTypeToString[v]
}

// MarshalJSON marshals the enum as a string
func (v ValidationType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
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
