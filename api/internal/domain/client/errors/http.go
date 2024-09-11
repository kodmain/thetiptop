package errors

const (
	// common errors
	ErrNoDto  = "no data transfer object provided"
	ErrNoData = "no data provided"

	// client errors
	ErrClientAlreadyExists    = "client already exists"
	ErrClientNotFound         = "client not found"
	ErrClientNotValidate      = "client not validate %s"
	ErrClientAlreadyValidated = "client already validated"

	// credential errors
	ErrCredentialNotFound      = "credential not found"
	ErrCredentialInvalid       = "credential invalid"
	ErrCredentialAlreadyExists = "credential already exists"

	// mail errors
	ErrMailSendFailed = "failed to send mail"

	// validation errors
	ErrValidationNotFound         = "validation not found"
	ErrValidationTokenNotFound    = "validation token not found"
	ErrValidationAlreadyValidated = "validation already validated"
	ErrValidationExpired          = "validation expired"

	// template
	ErrTemplateNotFound = "template %s not found"
)
