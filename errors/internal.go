package errors

import (
	"strings"
)

// Internal error codes
const (
	InternalErrorCode = "internal"
)

// InternalError holds errors which should not be exposed to clients
type InternalError struct {
	BaseError
	Cause error
}

// NewInternalError wraps original error with optional messages and gives InternalError
// Be careful when setting the msg, these messages will be visible to the client
func NewInternalError(cause error, msg ...string) InternalError {
	ierr := InternalError{
		BaseError: BaseError{
			Code:        InternalErrorCode,
			Message:     "Something went wrong",
			Description: cause.Error(),
		},
		Cause: cause,
	}
	if len(msg) > 0 {
		ierr.Message = strings.TrimSpace(strings.Join(msg, ", "))
	}
	return ierr
}
