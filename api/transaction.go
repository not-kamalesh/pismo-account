package api

import (
	"log/slog"
	"net/http"

	"github.com/not-kamalesh/pismo-account/dto"
)

func (api *APIHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	slog.Info("received request", "api", "CreateTransaction")

	ctx := r.Context()

	req, err := dto.ParseRequest(r, dto.NewCreateTransactionRequest)
	if err != nil {
		api.writeResponse(w, "CreateTransaction", nil, err, false)
		return
	}
	slog.Debug("parsed request", "api", "CreateTransaction", "request", req)

	resp, err := api.idempotencyMgr.Execute(req.ReferenceID, req.Hash(), func() (interface{}, error) {
		return api.transactionHandler.Create(ctx, req)
	})
	api.writeResponse(w, "CreateTransaction", resp, err, false)
}
