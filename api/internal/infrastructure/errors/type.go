package errors

type ErrorInterface interface {
	Error() string
	Code() int
}

type Error struct {
	code    int
	message string
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Error() string {
	return e.message
}

/*
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
	}
}
*/

func New(code int, message string) *Error {
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
