package errors

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound        = errors.New("record not found")
	ErrRecordNotFound  = gorm.ErrRecordNotFound
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidUserData = errors.New("invalid user data")
)

type Error struct {
	Code    int
	Message string
	Details error
}

type ErrorDTO struct {
	Code    int
	Message string
	Details string
}

func NewError(code int, message string, details error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func (e Error) Error() string {
	return e.Details.Error()
}

func (e Error) ToJson() *ErrorDTO {
	return &ErrorDTO{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details.Error(),
	}
}
