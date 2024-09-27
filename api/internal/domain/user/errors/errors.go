package errors

import (
	"net/http"
)

var (
	registredErrors = make(map[string]*Error)

	// Common errors
	ErrNoDto          = new(http.StatusBadRequest, "common.no_dto")
	ErrNoData         = new(http.StatusBadRequest, "common.no_data")
	ErrUnknown        = new(http.StatusInternalServerError, "common.unknown")
	ErrUnauthorized   = new(http.StatusUnauthorized, "common.unauthorized")
	ErrBadRequest     = new(http.StatusBadRequest, "common.bad_request")
	ErrInternalServer = new(http.StatusInternalServerError, "common.internal_error")

	// User errors
	ErrUserNotFound = new(http.StatusNotFound, "user.not_found")

	// Client errors
	ErrClientNotValidate      = new(http.StatusBadRequest, "client.not_validate")
	ErrClientNotFound         = new(http.StatusNotFound, "client.not_found")
	ErrClientAlreadyExists    = new(http.StatusConflict, "client.already_exists")
	ErrClientAlreadyValidated = new(http.StatusConflict, "client.already_validated")

	// Employee errors
	ErrEmployeeNotValidate      = new(http.StatusBadRequest, "employee.not_validate")
	ErrEmployeeNotFound         = new(http.StatusNotFound, "employee.not_found")
	ErrEmployeeAlreadyExists    = new(http.StatusConflict, "employee.already_exists")
	ErrEmployeeAlreadyValidated = new(http.StatusConflict, "employee.already_validated")

	// Credential errors
	ErrCredentialNotFound      = new(http.StatusNotFound, "credential.not_found")
	ErrCredentialAlreadyExists = new(http.StatusConflict, "credential.already_exists")

	// Mail errors
	ErrMailSendFailed = new(http.StatusInternalServerError, "mail.send_failed")

	// Validation errors
	ErrValidationNotFound         = new(http.StatusNotFound, "validation.not_found")
	ErrValidationTokenNotFound    = new(http.StatusNotFound, "validation.token_not_found")
	ErrValidationAlreadyValidated = new(http.StatusConflict, "validation.already_validated")
	ErrValidationExpired          = new(http.StatusGone, "validation.expired")

	// Template errors
	ErrTemplateNotFound = new(http.StatusNotFound, "template.not_found")
)
