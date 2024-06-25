package transfert

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

type Client struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewClient(obj data.Object, mandatory data.Validator) (*Client, error) {
	if err := mandatory.Check(obj); err != nil {
		return nil, err
	}

	return &Client{
		Email:    *obj.Get("email"),
		Password: *obj.Get("password"),
	}, nil
}
