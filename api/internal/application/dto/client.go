package dto

import (
	"errors"

	"github.com/kodmain/thetiptop/api/internal/application/validator"
)

type Client struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewClient(email, password string) (*Client, error) {
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
