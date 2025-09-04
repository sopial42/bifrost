package errors

import (
	"fmt"
	"log"
)

type AppErrorCode int

const (
	CodeErrUnknown AppErrorCode = iota
	CodeErrUnexpected
	CodeErrNotFound
	CodeErrUnauthorized
	CodeErrForbidden
	CodeErrInvalidInput
	CodeErrAlreadyExists
)

// Error list to use with errors.Is
var (
	ErrUnknown       = &AppError{Code: CodeErrUnknown}
	ErrUnexpected    = &AppError{Code: CodeErrUnexpected}
	ErrNotFound      = &AppError{Code: CodeErrNotFound}
	ErrUnauthorized  = &AppError{Code: CodeErrUnauthorized}
	ErrForbidden     = &AppError{Code: CodeErrForbidden}
	ErrInvalidInput  = &AppError{Code: CodeErrInvalidInput}
	ErrAlreadyExists = &AppError{Code: CodeErrAlreadyExists}
)

type AppError struct {
	Code    AppErrorCode `json:"app_code"`
	Message string       `json:"message"`
	Origin  error        `json:"-"`
}

func (e AppError) Error() string {
	if e.Origin != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Origin)
	}
	return e.Message
}

func (e AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		log.Printf("CanNOT cast target to AppError")
		return false
	}

	return e.Code == t.Code
}

func (e AppError) Unwrap() error {
	return e.Origin
}

func NewInvalidInput(message string, err error) *AppError {
	return &AppError{
		Code:    CodeErrInvalidInput,
		Message: message,
		Origin:  err,
	}
}

func NewUnauthorized(message string, err error) *AppError {
	return &AppError{
		Code:    CodeErrUnauthorized,
		Message: message,
		Origin:  err,
	}
}

func NewForbidden(message string, err error) *AppError {
	return &AppError{
		Code:    CodeErrForbidden,
		Message: message,
		Origin:  err,
	}
}

func NewNotFound(message string) *AppError {
	return &AppError{
		Code:    CodeErrNotFound,
		Message: message,
		Origin:  nil,
	}
}

func NewUnexpected(message string, err error) *AppError {
	return &AppError{
		Code:    CodeErrUnexpected,
		Message: message,
		Origin:  err,
	}
}

func NewAlreadyExists(message string) *AppError {
	return &AppError{
		Code:    CodeErrAlreadyExists,
		Message: message,
		Origin:  nil,
	}
}
