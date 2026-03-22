package api

import (
	"log/slog"
	"net/http"

	"github.com/not-kamalesh/pismo-account/dto"
)

func (api *APIHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	slog.Info("received request", "api", "CreateAccount")

	ctx := r.Context()

	req, err := dto.ParseRequest(r, dto.NewCreateAccountRequest)
	if err != nil {
		api.writeResponse(w, "CreateAccount", nil, err, false)
		return
	}

	resp, err := api.accountHandler.Create(ctx, req)

	api.writeResponse(w, "CreateAccount", resp, err, false)
}

func (api *APIHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	slog.Info("received request", "api", "GetAccount")

	ctx := r.Context()

	req, err := dto.ParseRequest(r, dto.NewGetAccountRequest)
	if err != nil {
		api.writeResponse(w, "GetAccount", nil, err, false)
		return
	}

	resp, err := api.accountHandler.Get(ctx, req)

	api.writeResponse(w, "GetAccount", resp, err, false)
}
