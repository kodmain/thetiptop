package transfert

import (
	"errors"

	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

type Client struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewClient(obj data.Object) (*Client, error) {
	email := obj.Get("email")
	password := obj.Get("password")

	err := errors.Join(
		validator.Email(email),
		validator.Password(password),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		Email:    email,
		Password: password,
	}, nil
}
