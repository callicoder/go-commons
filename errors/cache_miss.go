package errors

import (
	"strings"
)

// CacheMiss Error codes
const (
	CacheMissErrorCode = "cache_miss"
)

// CacheMissError holds cache miss errors
type CacheMissError struct {
	BaseError
	Cause error
}

// NewCacheMissError wraps original error with optional messages and gives CacheMissError
func NewCacheMissError(cause error, msg ...string) CacheMissError {
	cmerr := CacheMissError{
		BaseError: BaseError{
			Code:        CacheMissErrorCode,
			Message:     "Key not found in cache",
			Description: cause.Error(),
		},
		Cause: cause,
	}
	if len(msg) > 0 {
		cmerr.Message = strings.TrimSpace(strings.Join(msg, ", "))
	}
	return cmerr
}

func IsCacheMissError(err error) bool {
	_, ok := err.(CacheMissError)
	return ok
}
