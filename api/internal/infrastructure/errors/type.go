package errors

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
		message: e.message,
		data:    a,
	}
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
	err := &Error{
		code:    code,
		message: message,
	}

	registredErrors[message] = err

	return err
}

func ListErrors() map[string]*Error {
	return registredErrors
}
