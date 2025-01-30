package errors

import (
	"encoding/json"
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
)

type ErrorInterface interface {
	Error() string
	Code() int
	Log(error) *Error
}

type Errors map[string]ErrorInterface

func (es Errors) ToErrorInterface() ErrorInterface {
	if len(es) > 0 {
		return es
	}

	return nil
}

func (es Errors) Error() string {
	// Map to hold the error messages
	errorMap := make(map[string]string)
	for key, err := range es {
		errorMap[key] = err.Error()
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(errorMap)
	if err != nil {
		// Fallback in case of JSON marshalling error
		return fmt.Sprintf("failed to generate JSON: %v", err)
	}

	return string(jsonData)
}

func (es Errors) Code() int {
	for _, err := range es {
		return err.Code()
	}

	return 0
}

func (es Errors) Log(err error) *Error {
	for _, err := range es {
		err.Log(err)
	}

	return nil
}

func (es *Errors) Add(key string, err ErrorInterface) error {
	if _, ok := (*es)[key]; ok {
		return New(500, "Multiple errors")
	}

	(*es)[key] = err

	return nil
}

type Error struct {
	code    int
	message string
	logged  bool
}

func (e *Error) Log(err error) *Error {
	if !e.logged {
		e.logged = logger.Error(err)
	}

	return e
}

func (e Error) MarshalJSON() ([]byte, error) {
	exported := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    e.code,
		Message: e.message,
	}

	return json.Marshal(exported)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Error() string {
	return e.message
}

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
