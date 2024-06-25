package transfert

import (
	"errors"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
)

type Validation struct {
	Token    token.Luhn `json:"token"`
	ClientID string     `json:"client_id"`
}

func NewValidation(obj data.Object) (*Validation, error) {
	token := token.Luhn(obj.Get("token"))
	client := obj.Get("clientId")

	err := errors.Join(
		token.Validate(),
	)

	if err != nil {
		return nil, err
	}

	return &Validation{
		Token:    token,
		ClientID: client,
	}, nil
}
