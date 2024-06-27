package transfert

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
)

type Validation struct {
	Token    *token.Luhn `json:"token"`
	ClientID *string     `json:"client_id"`
}

func NewValidation(obj data.Object, mandatory data.Validator) (*Validation, error) {
	if err := mandatory.Check(obj); err != nil {
		return nil, err
	}

	luhn := obj.Get("token")
	client := obj.Get("clientId")
	token := token.NewLuhn(*luhn)

	return &Validation{
		Token:    &token,
		ClientID: client,
	}, nil
}
