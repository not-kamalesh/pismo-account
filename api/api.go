package api

import (
	"encoding/json"
	commonerrors "errors"
	"log/slog"
	"net/http"

	"github.com/not-kamalesh/pismo-account/errors"
	"github.com/not-kamalesh/pismo-account/internal/healthcheck"
)

type APIHandler struct {
	healthCheckHandler *healthcheck.Handler
}

func NewAPIHandler(healthcheckHandler *healthcheck.Handler) *APIHandler {
	apiHandler := &APIHandler{
		healthCheckHandler: healthcheckHandler,
	}
	return apiHandler
}

func (api *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	slog.Info("received request", "api", "HealthCheck")
	resp := api.healthCheckHandler.HealthCheck()
	api.writeResponse(w, "HealthCheck", resp, nil, false)
}

func (api *APIHandler) writeResponse(w http.ResponseWriter, apiName string, apiResp interface{}, apiErr error, isNilResponse bool) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case apiErr != nil:
		// Default to internal server error status code
		statusCode := http.StatusInternalServerError
		var pErr *errors.PismoError
		if commonerrors.As(apiErr, &pErr) {
			statusCode = pErr.GetHTTPCode()
			// Return only the PismoError struct as JSON, not the generic error interface
			apiErr = pErr
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
