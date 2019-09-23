package errors

// Validation error codes
const (
	CodeValidationError  = "validation_failed"
	CodeMissingField     = "missing_field"
	CodeMalformedRequest = "malformed_request"
)

type ValidationError struct {
	BaseError
	Details map[string]interface{}
}

// NewInvalidValueError should be used when a field in client request has
// invalid value. The field name passed will be added to details as 'parameter'
func NewInvalidValueError(field string, reason string) ValidationError {
	return ValidationError{
		BaseError: BaseError{
			Code:        CodeValidationError,
			Message:     "A parameter has invalid value",
			Description: "Validation Error",
		},
		Details: map[string]interface{}{
			"parameter": field,
			"reason":    reason,
		},
	}
}

// NewMissingFieldError should be used when a required field is missing or
// has an empty value in a client request.
func NewMissingFieldError(field string) ValidationError {
	return ValidationError{
		BaseError: BaseError{
			Code:        CodeMissingField,
			Message:     "At least one required parameter is missing",
			Description: "Validation Error",
		},
		Details: map[string]interface{}{
			"parameter": field,
		},
	}
}

// NewMalformedRequestError should be used when request data is malformed.
func NewMalformedRequestError(reason string) ValidationError {
	return ValidationError{
		BaseError: BaseError{
			Code:        CodeMalformedRequest,
			Message:     "Request payload is invalid or malformed",
			Description: "Malformed Request",
		},
		Details: map[string]interface{}{
			"reason": reason,
		},
	}
}

// NewValidationError should be used to generate custom validation errors
func NewValidationError(msg string) ValidationError {
	return ValidationError{
		BaseError: BaseError{
			Code:        CodeValidationError,
			Message:     msg,
			Description: "Validation Error",
		},
	}
}

func (err ValidationError) ErrorCode() string {
	if _, found := err.Details["parameter"]; found {
		return err.Code + ":" + err.Details["parameter"].(string)
	}
	return err.Code
}
