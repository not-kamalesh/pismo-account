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
	// Bad Request Errors - Validation
	ErrInvalidArgument      = New(PismoErrorCodeInvalidArgument, "Invalid argument provided", http.StatusBadRequest)
	ErrInvalidMsgID         = New(PismoErrorCodeInvalidArgument, "Invalid msg_id", http.StatusBadRequest)
	ErrInvalidReferenceID   = New(PismoErrorCodeInvalidArgument, "Invalid reference_id", http.StatusBadRequest)
	ErrInvalidDocumentID    = New(PismoErrorCodeInvalidArgument, "Invalid document_id", http.StatusBadRequest)
	ErrInvalidAccountID     = New(PismoErrorCodeInvalidArgument, "Invalid account_id", http.StatusBadRequest)
	ErrInvalidCurrency      = New(PismoErrorCodeInvalidArgument, "Invalid currency", http.StatusBadRequest)
	ErrInvalidAmount        = New(PismoErrorCodeInvalidArgument, "Invalid amount, amount should be positive", http.StatusBadRequest)
	ErrInvalidOperationType = New(PismoErrorCodeInvalidArgument, "Invalid operation_type_id", http.StatusBadRequest)

	ErrNotFound         = New(PismoErrorCodeNotFound, "Requested resource not found", http.StatusNotFound)
	ErrAlreadyExists    = New(PismoErrorCodeAlreadyExists, "Resource already exists", http.StatusConflict)
	ErrPermissionDenied = New(PismoErrorCodePermissionDenied, "Permission denied", http.StatusForbidden)
	ErrInternal         = New(PismoErrorCodeInternal, "Internal server error", http.StatusInternalServerError)
)
