package apperrors

import "net/http"

type ErrorType string

const (
	TypeNotFound     ErrorType = "NOT_FOUND"
	TypeInvalidInput ErrorType = "INVALID_INPUT"
	TypeConflict     ErrorType = "CONFLICT"
	TypeInternal     ErrorType = "INTERNAL"
)

type AppError struct {
	StatusCode int
	Type       ErrorType
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NotFoundError(msg string, err error) *AppError {
	return &AppError{StatusCode: http.StatusNotFound, Type: TypeNotFound, Message: msg, Err: err}
}

func InvalidInputError(msg string, err error) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Type: TypeInvalidInput, Message: msg, Err: err}
}

func ConflictError(msg string, err error) *AppError {
	return &AppError{StatusCode: http.StatusConflict, Type: TypeConflict, Message: msg, Err: err}
}

func InternalError(msg string, err error) *AppError {
	return &AppError{StatusCode: http.StatusInternalServerError, Type: TypeInternal, Message: msg, Err: err}
}
