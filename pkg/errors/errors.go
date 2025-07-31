package errors

import (
	"errors"
	"fmt"
)

type AppErrorCode int

const (
	ErrUnexpected AppErrorCode = iota + 1
	ErrNotFound
	ErrUnauthorized
	ErrForbidden
	ErrInvalidInput
	ErrAlreadyExists
)

type AppError struct {
	Code    AppErrorCode `json:"code"`
	Message string       `json:"message"`
	Origin  error        `json:"-"`
}

func (e AppError) Error() string {
	if e.Origin != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Origin)
	}
	return e.Message
}

func (e AppError) Unwrap() error {
	return e.Origin
}

func NewInvalidInput(message string, err error) *AppError {
	return &AppError{
		Code:    ErrInvalidInput,
		Message: message,
		Origin:  err,
	}
}

func NewUnauthorized(message string, err error) *AppError {
	return &AppError{
		Code:    ErrUnauthorized,
		Message: message,
		Origin:  err,
	}
}

func NewForbidden(message string, err error) *AppError {
	return &AppError{
		Code:    ErrForbidden,
		Message: message,
		Origin:  err,
	}
}

func NewNotFound(message string) *AppError {
	return &AppError{
		Code:    ErrNotFound,
		Message: message,
		Origin:  nil,
	}
}

func AsNotFound(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	for {
		if errors.As(err, &appErr) {
			if appErr.Code == ErrNotFound {
				return true
			}
		}

		err = errors.Unwrap(err)
		if err == nil {
			break
		}
	}

	return false
}

func AsAlreadyExists(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	for {
		if errors.As(err, &appErr) {
			if appErr.Code == ErrAlreadyExists {
				return true
			}
		}

		err = errors.Unwrap(err)
		if err == nil {
			break
		}
	}

	return false
}

func NewUnexpected(message string, err error) *AppError {
	return &AppError{
		Code:    ErrUnexpected,
		Message: message,
		Origin:  err,
	}
}

func NewAlreadyExists(message string) *AppError {
	return &AppError{
		Code:    ErrAlreadyExists,
		Message: message,
		Origin:  nil,
	}
}
