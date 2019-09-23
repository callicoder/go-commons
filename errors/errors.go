package errors

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

// New returns an error with supplied message
func New(msg string) error {
	return pkgerrors.New(msg)
}

// Wrap returns new error by annotating the passed error with the message
func Wrap(err error, msg string) error {
	return pkgerrors.Wrap(err, msg)
}

// Wrapf returns a new error by annotating passed error with formated message
func Wrapf(err error, msg string, args ...interface{}) error {
	return pkgerrors.Wrapf(err, msg, args...)
}

// Cause returns the underlying cause of the error by unwrapping the error
// StackTrace.
func Cause(err error) error {
	return pkgerrors.Cause(err)
}

// BaseError holds code and description of an error
type BaseError struct {
	//Unique Code for an error
	Code string `json:"code"`
	//Message is the one that is sent to the client
	Message string `json:"message"`
	//Description is the long description that gives further detail about the error. This is not sent to the client
	Description string `json:"-"`
}

func (err BaseError) Error() string {
	message := fmt.Sprintf("[%s] %s", err.Code, err.Message)
	if err.Description != "" {
		message = message + " => " + err.Description
	}
	return message
}

// ErrorCode returns the error code
func (err BaseError) ErrorCode() string {
	return err.Code
}
