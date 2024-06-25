package errors

const (
	// common errors
	ErrNoDto = "no data transfer object provided"

	// client errors
	ErrClientAlreadyExists    = "client already exists"
	ErrClientNotFound         = "client not found"
	ErrClientNotValidated     = "client not validated"
	ErrClientAlreadyValidated = "client already validated"

	// mail errors
	ErrMailSendFailed = "failed to send mail"

	// validation errors
	ErrValidationNotFound         = "validation not found"
	ErrValidationAlreadyValidated = "validation already validated"
	ErrValidationExpired          = "validation expired"
)
