package validator

import (
	"errors"
	"net/mail"
	"strings"
	"unicode"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
)

const (
	CAN_BE_NIL  = true
	CANT_BE_NIL = false
)

func Required(value *string) error {
	if value == nil {
		return errors.New("value is required")
	}

	return nil
}

func Email(email *string) error {
	_, err := mail.ParseAddress(*email)
	return err
}

func Luhn(value *string) error {
	if value == nil {
		return errors.New("value is required")
	}

	luhn := token.Luhn(*value)
	return luhn.Validate()
}

func ID(uuid *string) error {
	id := *uuid
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return errors.New("invalid UUID")
	}

	return nil
}

func Password(password *string) error {
	var (
		hasMin     bool
		hasMaj     bool
		hasNumber  bool
		hasSpecial bool
		errs       []string
	)

	if len(*password) < 8 {
		errs = append(errs, "- password is too short")
	}

	if len(*password) > 64 {
		errs = append(errs, "- password is too long")
	}

	for _, c := range *password {
		switch {
		case unicode.IsLower(c):
			hasMin = true
		case unicode.IsUpper(c):
			hasMaj = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasMin {
		errs = append(errs, "- password must include lowercase letters")
	}

	if !hasMaj {
		errs = append(errs, "- password must include uppercase letters")
	}

	if !hasNumber {
		errs = append(errs, "- password must include numbers")
	}

	if !hasSpecial {
		errs = append(errs, "- password must include special characters")
	}

	if len(errs) > 0 {
		return errors.New("invalid password: \n" + strings.Join(errs, "\n"))
	}

	return nil
}
