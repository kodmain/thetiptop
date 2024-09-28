package errors

import (
	"encoding/json"
)

type ErrorInterface interface {
	Error() string
	Code() int
	Data() any
	WithData(a ...any) *Error
}

type Error struct {
	code    int
	message string
	data    any
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Data() any {
	return e.data
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) WithData(a ...any) *Error {
	return &Error{
		code:    e.code,
		message: e.Error(),
		data:    a,
	}
}

func (e *Error) Marshal() ([]byte, error) {
	serialized := map[string]any{
		"code":    e.code,
		"message": e.Error(),
		"data":    e.data,
	}
	return json.Marshal(serialized)
}

func FromErr(err error, code ErrorInterface) ErrorInterface {
	if err == nil {
		return nil
	}

	if e, ok := err.(ErrorInterface); ok {
		return e
	}

	return &Error{
		code:    code.Code(),
		message: err.Error(),
		data:    nil,
	}
}

func new(code int, message string) *Error {
	registredErrors[message] = &Error{
		code:    code,
		message: message,
		data:    nil,
	}

	return registredErrors[message]
}
