package errors

import (
	"net/http"
)

var (
	registredErrors = make(map[string]*Error)

	// Common errors
	ErrNoDto          = New(http.StatusBadRequest, "common.no_dto")
	ErrNoData         = New(http.StatusBadRequest, "common.no_data")
	ErrUnauthorized   = New(http.StatusUnauthorized, "common.unauthorized")
	ErrBadRequest     = New(http.StatusBadRequest, "common.bad_request")
	ErrForbidden      = New(http.StatusForbidden, "common.forbidden")
	ErrNotFound       = New(http.StatusNotFound, "common.not_found")
	ErrInternalServer = New(http.StatusInternalServerError, "common.internal_error")

	// Hash errors
	ErrHashAlgoUnknown = New(http.StatusInternalServerError, "hash.algo_unknown")
	ErrHashMismatch    = New(http.StatusUnauthorized, "hash.mismatch")

	// Validator errors
	ErrValueRequired                     = New(http.StatusBadRequest, "validator.required")
	ErrValueIsNotString                  = New(http.StatusBadRequest, "validator.is_not_string")
	ErrValueIsNotInt                     = New(http.StatusBadRequest, "validator.is_not_int")
	ErrValueIsNotNumber                  = New(http.StatusBadRequest, "validator.is_not_number")
	ErrValueIsNotBool                    = New(http.StatusBadRequest, "validator.is_not_bool")
	ErrValueBoolMustBeTrue               = New(http.StatusBadRequest, "validator.bool_must_be_true")
	ErrValueBoolMustBeFalse              = New(http.StatusBadRequest, "validator.bool_must_be_false")
	ErrValueIsNotFloat                   = New(http.StatusBadRequest, "validator.is_not_float")
	ErrValueIsNotEmail                   = New(http.StatusBadRequest, "validator.is_not_email")
	ErrValueIsNotPassword                = New(http.StatusBadRequest, "validator.is_not_password")
	ErrValuePasswordIsToShort            = New(http.StatusBadRequest, "validator.password_is_to_short")
	ErrValuePasswordIsToLong             = New(http.StatusBadRequest, "validator.password_is_to_long")
	ErrValuePasswordMustIncludeLowercase = New(http.StatusBadRequest, "validator.password_must_include_lowercase")
	ErrValuePasswordMustIncludeUppercase = New(http.StatusBadRequest, "validator.password_must_include_uppercase")
	ErrValuePasswordMustIncludeNumber    = New(http.StatusBadRequest, "validator.password_must_include_number")
	ErrValuePasswordMustIncludeSpecial   = New(http.StatusBadRequest, "validator.password_must_include_special")
	ErrValueIsNotPhone                   = New(http.StatusBadRequest, "validator.is_not_phone")
	ErrValueIsNotID                      = New(http.StatusBadRequest, "validator.is_not_id")
	ErrValueIsNotLuhn                    = New(http.StatusBadRequest, "validator.is_not_luhn")
	ErrValueIsNotURL                     = New(http.StatusBadRequest, "validator.is_not_url")
	ErrValueIsNotDate                    = New(http.StatusBadRequest, "validator.is_not_date")
	ErrValueIsNotTime                    = New(http.StatusBadRequest, "validator.is_not_time")
	ErrValueIsNotUUID                    = New(http.StatusBadRequest, "validator.is_not_uuid")

	// Auth errors
	ErrAuthInvalidToken = New(http.StatusBadRequest, "auth.invalid_token")
	ErrAuthFailed       = New(http.StatusUnauthorized, "auth.failed")
	ErrAuthBadFormat    = New(http.StatusBadRequest, "auth.bad_format")
	ErrAuthForbidden    = New(http.StatusForbidden, "auth.forbidden")

	// Mail errors
	ErrMailSendFailed = New(http.StatusInternalServerError, "mail.send_failed")

	// Template errors
	ErrMailTemplateNotFound = New(http.StatusNotFound, "template.mail.not_found")
)
