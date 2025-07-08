package entity

import (
	"fmt"
	"net/http"
)

type Error struct {
	Message string
	Detail  string
	Code    int
	Err     error
}

func NewError(message string, code int) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}

func (e *Error) WithError(err error) *Error {
	if err != nil {
		obj := *e
		obj.Err = err
		obj.Detail = err.Error()

		return &obj
	}
	return e
}

func (e *Error) Decorate(msg string) *Error {
	obj := *e
	obj.Detail = msg
	return &obj
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}

	return e.Message
}

var (
	NotFoundError       = NewError("not found", http.StatusNotFound)
	InternalServerError = NewError("internal server error", http.StatusInternalServerError)
	DuplicateKeyError   = NewError("duplicate key error", http.StatusConflict)

	InvalidBodyError = NewError("invalid request body", http.StatusBadRequest)
)
