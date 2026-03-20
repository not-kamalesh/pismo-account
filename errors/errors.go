package errors

import (
	"fmt"
	"net/http"
)

type PismoError struct {
	Code     PismoErrorCode `json:"code"`
	Message  string         `json:"message"`
	Err      error          `json:"-"`
	HTTPCode int            `json:"-"`
}

func New(code PismoErrorCode, msg string, httpCode int) *PismoError {
	return &PismoError{
		Code:     code,
		Message:  msg,
		HTTPCode: httpCode,
	}
}

func (e *PismoError) Error() string {
	return fmt.Sprintf("code: %s, message: %s, details: %s", e.Code, e.Message, e.Err.Error())
}

func (e *PismoError) GetCode() PismoErrorCode {
	return e.Code
}

func (e *PismoError) GetHTTPCode() int {
	if e.HTTPCode != 0 {
		return e.HTTPCode
	}
	return http.StatusInternalServerError
}

// Pre-defined errors @TODO to add relevent errors and error codes
var (
	ErrInvalidArgument  = New(PismoErrorCodeInvalidArgument, "Invalid argument provided", http.StatusBadRequest)
	ErrNotFound         = New(PismoErrorCodeNotFound, "Requested resource not found", http.StatusNotFound)
	ErrAlreadyExists    = New(PismoErrorCodeAlreadyExists, "Resource already exists", http.StatusConflict)
	ErrPermissionDenied = New(PismoErrorCodePermissionDenied, "Permission denied", http.StatusForbidden)
	ErrInternal         = New(PismoErrorCodeInternal, "Internal server error", http.StatusInternalServerError)
)
