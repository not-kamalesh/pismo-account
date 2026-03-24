package api

import (
	"encoding/json"
	commonerrors "errors"
	"log/slog"
	"net/http"

	"github.com/not-kamalesh/pismo-account/errors"
	"github.com/not-kamalesh/pismo-account/internal/account"
	"github.com/not-kamalesh/pismo-account/internal/healthcheck"
	"github.com/not-kamalesh/pismo-account/internal/idempotencymgr"
	"github.com/not-kamalesh/pismo-account/internal/transaction"
)

type APIHandler struct {
	healthCheckHandler healthcheck.HealthCheckHandler
	accountHandler     account.AccountHandler
	transactionHandler transaction.TransactionHandler
	idempotencyMgr     idempotencymgr.IdempotencyMgr
}

func NewAPIHandler(healthcheckHandler healthcheck.HealthCheckHandler,
	accountHandler account.AccountHandler,
	transactionHandler transaction.TransactionHandler,
	idempotencyMgr idempotencymgr.IdempotencyMgr,
) *APIHandler {

	apiHandler := &APIHandler{
		healthCheckHandler: healthcheckHandler,
		accountHandler:     accountHandler,
		transactionHandler: transactionHandler,
		idempotencyMgr:     idempotencyMgr,
	}
	return apiHandler
}

// writeResponse is a helper function which writes the response and error to the response writer
func (api *APIHandler) writeResponse(w http.ResponseWriter, apiName string, apiResp interface{}, apiErr error, isNilResponse bool) {
	w.Header().Set("Content-Type", "application/json")
	slog.Debug("API Response", "apiName", apiName, "resp", apiResp, "error", apiErr)
	switch {
	case apiErr == nil && apiResp == nil && isNilResponse == false:
		// This case is not expected, when it happens, return internal server serror with no body
		w.WriteHeader(http.StatusInternalServerError)
	case apiErr != nil:
		// Default to internal server error status code
		statusCode := http.StatusInternalServerError
		var pErr *errors.PismoError
		if commonerrors.As(apiErr, &pErr) {
			statusCode = pErr.GetHTTPCode()
			// Return only the PismoError struct as JSON, not the generic error interface
			apiErr = pErr
		} else {
			apiErr = &errors.PismoError{
				Code:    errors.PismoErrorCodeInternal,
				Message: apiErr.Error(),
			}
		}
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(apiErr); err != nil {
			slog.Error("failed to write error response", "api", apiName, "error", err)
		}
	case isNilResponse:
		// No response to write, send 204
		w.WriteHeader(http.StatusNoContent)
	default:
		// Happy path, return response object
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(apiResp); err != nil {
			slog.Error("failed to write api response", "api", apiName, "resp", apiResp, "error", err)
		}
	}
}
