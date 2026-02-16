package errors

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error ErrorBody `json:"message"`
}

type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HTTPError struct {
	Status  int
	Code    int
	Message string

	e error
}

func (e *HTTPError) Error() string {
	if e.e != nil {
		return fmt.Sprintf("HTTP %d: %s - %v", e.Status, e.Message, e.e)
	}
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.e
}

func BadRequest(code int, message string, err error) *HTTPError {
	return &HTTPError{
		Status:  http.StatusBadRequest,
		Code:    code,
		Message: message,
		e:       err,
	}
}

func InternalServerError(code int, message string, err error) *HTTPError {
	return &HTTPError{
		Status:  http.StatusInternalServerError,
		Code:    code,
		Message: message,
		e:       err,
	}
}

func NotFound(code int, message string, err error) *HTTPError {
	return &HTTPError{
		Status:  http.StatusNotFound,
		Code:    code,
		Message: message,
		e:       err,
	}
}

func Unauthorized(code int, message string, err error) *HTTPError {
	return &HTTPError{
		Status:  http.StatusUnauthorized,
		Code:    code,
		Message: message,
		e:       err,
	}
}
