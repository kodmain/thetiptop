package errors_domain_user

import (
	"net/http"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

var (

	// User errors
	ErrUserNotFound = errors.New(http.StatusNotFound, "user.not_found")

	// Client errors
	ErrClientNotValidate      = errors.New(http.StatusBadRequest, "client.not_validate")
	ErrClientNotFound         = errors.New(http.StatusNotFound, "client.not_found")
	ErrClientAlreadyExists    = errors.New(http.StatusConflict, "client.already_exists")
	ErrClientAlreadyValidated = errors.New(http.StatusConflict, "client.already_validated")

	// Employee errors
	ErrEmployeeNotValidate      = errors.New(http.StatusBadRequest, "employee.not_validate")
	ErrEmployeeNotFound         = errors.New(http.StatusNotFound, "employee.not_found")
	ErrEmployeeAlreadyExists    = errors.New(http.StatusConflict, "employee.already_exists")
	ErrEmployeeAlreadyValidated = errors.New(http.StatusConflict, "employee.already_validated")

	// Credential errors
	ErrCredentialNotFound      = errors.New(http.StatusNotFound, "credential.not_found")
	ErrCredentialNotValid      = errors.New(http.StatusBadRequest, "credential.not_valid")
	ErrCredentialAlreadyExists = errors.New(http.StatusConflict, "credential.already_exists")

	// Validation errors
	ErrValidationNotFound         = errors.New(http.StatusNotFound, "validation.not_found")
	ErrValidationTokenNotFound    = errors.New(http.StatusNotFound, "validation.token_not_found")
	ErrValidationAlreadyValidated = errors.New(http.StatusConflict, "validation.already_validated")
	ErrValidationExpired          = errors.New(http.StatusGone, "validation.expired")
)
