package errors

import (
	"strings"
)

// NotFoundError codes
const (
	UnauthorizedErrorCode = "unauthorized"
)

// UnauthorizedError holds unauthorized error
type UnauthorizedError struct {
	BaseError
	Cause error
}

// NewUnauthorizedError wraps original error with optional messages and gives UnauthorizedError
func NewUnauthorizedError(cause error, msg ...string) UnauthorizedError {
	nferr := UnauthorizedError{
		BaseError: BaseError{
			Code:        UnauthorizedErrorCode,
			Message:     "Not authorized",
			Description: cause.Error(),
		},
		Cause: cause,
	}
	if len(msg) > 0 {
		nferr.Message = strings.TrimSpace(strings.Join(msg, ", "))
	}
	return nferr
}
