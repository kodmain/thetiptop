package validator

import (
	"net/mail"
	"reflect"
	"unicode"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
)

const (
	CAN_BE_NIL  = true
	CANT_BE_NIL = false
)

func Required(value any, name string) errors.ErrorInterface {
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		return errors.ErrValueRequired
	}

	return nil
}

func Email(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrString(value)
	if str == nil {
		return errors.ErrValueIsNotString
	}

	_, err := mail.ParseAddress(*str)
	if err != nil {
		return errors.ErrValueIsNotEmail
	}

	return nil
}

func Luhn(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrString(value)
	if str == nil {
		return errors.ErrValueIsNotString
	}

	luhn := token.Luhn(*str)

	return luhn.Validate()
}

func ID(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrString(value)
	if str == nil {
		return errors.ErrValueIsNotString
	}

	id := *str
	if len(id) != 36 || id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return errors.ErrValueIsNotUUID
	}

	return nil
}

func Password(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrString(value)
	if str == nil {
		return errors.ErrValueIsNotString
	}

	var (
		hasMin     bool
		hasMaj     bool
		hasNumber  bool
		hasSpecial bool
	)

	if len(*str) < 8 {
		return errors.ErrValuePasswordIsToShort
	}

	if len(*str) > 64 {
		return errors.ErrValuePasswordIsToLong
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
		return errors.ErrValuePasswordMustIncludeLowercase
	}

	if !hasMaj {
		return errors.ErrValuePasswordMustIncludeUppercase
	}

	if !hasNumber {
		return errors.ErrValuePasswordMustIncludeNumber
	}

	if !hasSpecial {
		return errors.ErrValuePasswordMustIncludeSpecial
	}

	return nil
}

func IsTrue(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrBool(value)
	if str == nil {
		return errors.ErrValueIsNotBool
	}

	if !*str {
		return errors.ErrValueBoolMustBeTrue
	}

	return nil
}

func IsFalse(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrBool(value)
	if str == nil {
		return errors.ErrValueIsNotBool
	}

	if *str {
		return errors.ErrValueBoolMustBeFalse
	}

	return nil
}

func IsBool(value any, name string) errors.ErrorInterface {
	if err := Required(value, name); err != nil {
		return err
	}

	str := anyToPtrBool(value)
	if str == nil {
		return errors.ErrValueIsNotBool
	}

	return nil
}
