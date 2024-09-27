package validator

import (
	"net/mail"
	"reflect"
	"strings"
	"unicode"

	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
)

const (
	CAN_BE_NIL  = true
	CANT_BE_NIL = false
)

func Required(value any, name string) errors.ErrorInterface {
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		return errors.ErrBadRequest.WithData("value %v is required", name)
	}

	return nil
}

func Email(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrString(value)
	if str == nil {
		return errors.ErrBadRequest.WithData("%v is not a string", name)
	}

	_, err := mail.ParseAddress(*str)
	return errors.FromErr(err, errors.ErrBadRequest)
}

func Luhn(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrString(value)
	if str == nil {
		return errors.ErrBadRequest.WithData("%v is not a string", name)
	}

	luhn := token.Luhn(*str)

	return errors.FromErr(luhn.Validate(), errors.ErrBadRequest)
}

func ID(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrString(value)
	if str == nil {
		return errors.ErrBadRequest.WithData("%v is not a string", name)
	}

	id := *str
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return errors.ErrBadRequest.WithData("invalid UUID")
	}

	return nil
}

func Password(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrString(value)
	if str == nil {
		return errors.ErrBadRequest.WithData("%v is not a string", name)
	}

	var (
		hasMin     bool
		hasMaj     bool
		hasNumber  bool
		hasSpecial bool
		errs       []string
	)

	if len(*str) < 8 {
		errs = append(errs, "- password is too short")
	}

	if len(*str) > 64 {
		errs = append(errs, "- password is too long")
	}

	for _, c := range *str {
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
		return errors.ErrBadRequest.WithData("invalid password: \n" + strings.Join(errs, "\n"))
	}

	return nil
}

func IsTrue(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrBool(value)
	if str == nil {
		return errors.ErrBadRequest.WithData("%v is not a boolean", name)
	}

	if !*str {
		return errors.ErrBadRequest.WithData("%v sould be true", name)
	}

	return nil
}

func IsFalse(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrBool(value)
	if str == nil {
		return errors.ErrBadRequest.WithData("%v is not a boolean", name)
	}

	if *str {
		return errors.ErrBadRequest.WithData("%v sould be false", name)
	}

	return nil
}

func IsBool(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrBool(value)
	if str == nil {
		return errors.ErrBadRequest.WithData("%v is not a boolean", name)
	}

	return nil
}
