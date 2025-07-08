package entity

import "net/http"

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
	e.Err = err
	return e
}

func (e *Error) Decorate(msg string) *Error {
	e.Detail = msg
	return e
}

func (e *Error) Error() string {
	//fmt.Println(e.Message)
	//if e.Detail != "" {
	//	return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	//}

	return e.Message
}

var (
	NotFoundError       = NewError("not found", http.StatusNotFound)
	InternalServerError = NewError("internal server error", http.StatusInternalServerError)
	DuplicateKeyError   = NewError("duplicate key error", http.StatusConflict)

	InvalidBodyError = NewError("invalid request body", http.StatusBadRequest)
)
