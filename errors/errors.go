package errors

import (
	"fmt"
	"net/http"

	"github.com/callicoder/go-commons/errors/codes"
	pkgerrors "github.com/pkg/errors"
)

func New(msg string) error {
	return &BaseError{
		Code:    codes.Internal,
		Message: msg,
		Stack:   pkgerrors.New(msg),
	}
}

func Newf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return &BaseError{
		Code:    codes.Internal,
		Message: msg,
		Stack:   pkgerrors.New(msg),
	}
}

func Wrap(err error, msg string) error {
	return &BaseError{
		Code:    codes.Internal,
		Message: msg,
		Stack:   pkgerrors.Wrap(err, msg),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return &BaseError{
		Code:    codes.Internal,
		Message: msg,
		Stack:   pkgerrors.Wrapf(err, msg),
	}
}

func WithDetails(details ...Detail) *withEntry {
	return &withEntry{
		details: details,
	}
}

func WithCode(code string) *withEntry {
	return &withEntry{
		code: code,
	}
}

func Cause(err error) error {
	return pkgerrors.Cause(err)
}

func HTTPStatus(err error) int64 {
	baseErr, ok := err.(*BaseError)
	if !ok {
		return http.StatusInternalServerError
	}

	return codes.HttpStatus(baseErr.Code)
}

type Detail struct {
	// Resource which has the error
	Resource string `json:"resource,omitempty"`
	// The specific field of the resource which has the error
	Field string `json:"field,omitempty"`
	// The value of the field which was erroneous
	Value interface{} `json:"value,omitempty"`
	// Message for this error
	Message string `json:"message,omitempty"`
}

// BaseError holds code, message, and details of an error
type BaseError struct {
	//Unique Code for an error
	Code string `json:"code"`
	//Message is the one that is sent to the client
	Message string `json:"message"`
	//Details is any additional details related to the error
	Details []Detail `json:"details,omitempty"`
	//Stack is used only for developers and not exposed in the json serialization
	Stack error `json:"-"`
}

func (err *BaseError) Error() string {
	return err.Stack.Error()
}

func (err *BaseError) Cause() error {
	return pkgerrors.Cause(err.Stack)
}

type withEntry struct {
	code    string
	details []Detail
}

func (entry *withEntry) New(msg string) error {
	return &BaseError{
		Code:    entry.code,
		Message: msg,
		Details: entry.details,
		Stack:   pkgerrors.New(msg),
	}
}

func (entry *withEntry) Newf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return &BaseError{
		Code:    entry.code,
		Message: msg,
		Details: entry.details,
		Stack:   pkgerrors.New(msg),
	}
}

func (entry *withEntry) Wrap(err error, msg string) error {
	return &BaseError{
		Code:    entry.code,
		Message: msg,
		Details: entry.details,
		Stack:   pkgerrors.Wrap(err, msg),
	}
}

func (entry *withEntry) Wrapf(err error, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return &BaseError{
		Code:    entry.code,
		Message: msg,
		Details: entry.details,
		Stack:   pkgerrors.Wrapf(err, msg),
	}
}

func (entry *withEntry) WithCode(code string) *withEntry {
	return &withEntry{
		details: entry.details,
		code:    code,
	}
}

func (entry *withEntry) WithDetails(details ...Detail) *withEntry {
	return &withEntry{
		details: append(entry.details, details...),
		code:    entry.code,
	}
}
