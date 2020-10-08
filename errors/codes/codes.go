package codes

import "net/http"

// string error codes allow us to incorporate application specific descriptive error codes too.
// otherwise, an int type for code would have been better
const (
	BadRequest   = "bad_request"
	Unauthorized = "unauthorized"
	Forbidden    = "forbidden"
	NotFound     = "not_found"
	Conflict     = "conflict"
	Internal     = "internal"
)

var codeToHttpStatus = map[string]int64{
	BadRequest:   http.StatusBadRequest,
	Unauthorized: http.StatusUnauthorized,
	Forbidden:    http.StatusForbidden,
	NotFound:     http.StatusNotFound,
	Conflict:     http.StatusConflict,
}

func HttpStatus(code string) int64 {
	if status, ok := codeToHttpStatus[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}
