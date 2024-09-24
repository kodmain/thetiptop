package errors

import "github.com/gofiber/fiber/v2"

func HTTPResponse(err error) (int, string) {
	switch err.Error() {
	case ErrClientNotFound:
	case ErrEmployeeNotFound:
		return fiber.StatusNotFound, err.Error()
	case ErrClientAlreadyExists:
		return fiber.StatusConflict, err.Error()
	case ErrUnauthorized:
		return fiber.StatusUnauthorized, err.Error()
	}
	return fiber.StatusInternalServerError, err.Error()
}
