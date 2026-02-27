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
	ip      string

	e error
}

func (e *HTTPError) Error() string {
	if e.e != nil {
		return fmt.Sprintf("%s HTTP %d: %s - %v", e.ip, e.Status, e.Message, e.e)
	}
	return fmt.Sprintf("%s HTTP %d: %s", e.ip, e.Status, e.Message)
}

func (e *HTTPError) Unwrap() error {
	return e.e
}

func BadRequest(code int, message string, err error, r *http.Request) *HTTPError {
	return &HTTPError{
		Status:  http.StatusBadRequest,
		Code:    code,
		Message: message,
		e:       err,
		ip:      r.RemoteAddr,
	}
}

func InternalServerError(code int, message string, err error, r *http.Request) *HTTPError {
	return &HTTPError{
		Status:  http.StatusInternalServerError,
		Code:    code,
		Message: message,
		e:       err,
		ip:      r.RemoteAddr,
	}
}

func NotFound(code int, message string, err error, r *http.Request) *HTTPError {
	return &HTTPError{
		Status:  http.StatusNotFound,
		Code:    code,
		Message: message,
		e:       err,
		ip:      r.RemoteAddr,
	}
}

func Unauthorized(code int, message string, err error, r *http.Request) *HTTPError {
	return &HTTPError{
		Status:  http.StatusUnauthorized,
		Code:    code,
		Message: message,
		e:       err,
		ip:      r.RemoteAddr,
	}
}

func Forbidden(code int, message string, err error, r *http.Request) *HTTPError {
	return &HTTPError{
		Status:  http.StatusForbidden,
		Code:    code,
		Message: message,
		e:       err,
		ip:      r.RemoteAddr,
	}
}

func MethodNotAllowed(code int, message string, err error, r *http.Request) *HTTPError {
	return &HTTPError{
		Status:  http.StatusMethodNotAllowed,
		Code:    code,
		Message: message,
		e:       err,
		ip:      r.RemoteAddr,
	}
}
