package errors

import (
	"strings"
)

// NotFoundError codes
const (
	BadRequestErrorCode = "bad_request"
)

// UnauthorizedError holds unauthorized error
type BadRequestError struct {
	BaseError
	Cause error
}

// NewUnauthorizedError wraps original error with optional messages and gives UnauthorizedError
func NewBadRequestError(cause error, msg ...string) BadRequestError {
	berr := BadRequestError{
		BaseError: BaseError{
			Code:        BadRequestErrorCode,
			Message:     "Bad Request",
			Description: cause.Error(),
		},
		Cause: cause,
	}
	if len(msg) > 0 {
		berr.Message = strings.TrimSpace(strings.Join(msg, ", "))
	}
	return berr
}
