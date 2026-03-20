package api

import (
	"log/slog"
	"net/http"
)

func (api *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	slog.Info("received request", "api", "HealthCheck")

	resp := api.healthCheckHandler.HealthCheck()

	api.writeResponse(w, "HealthCheck", resp, nil, false)
}
