package dto

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/not-kamalesh/pismo-account/common/types"
	"github.com/not-kamalesh/pismo-account/errors"
)

type CreateTransactionRequest struct {
	MsgID           string              `json:"msg_id"`
	ReferenceID     string              `json:"reference_id"`
	AccountID       int64               `json:"account_id"`
	OperationTypeID types.OperationType `json:"operation_type_id"`
	Amount          float64             `json:"amount"`
}

func NewCreateTransactionRequest() *CreateTransactionRequest {
	return &CreateTransactionRequest{}
}

func (r *CreateTransactionRequest) Parse(httpReq *http.Request) error {
	decoder := json.NewDecoder(httpReq.Body)
	if err := decoder.Decode(r); err != nil {
		slog.Error("error occured to decode", "err", err)
		return err
	}

	return nil
}

func (r *CreateTransactionRequest) Validate() error {
	if r.MsgID == "" {
		return errors.ErrInvalidMsgID
	}
	if r.ReferenceID == "" {
		return errors.ErrInvalidReferenceID
	}
	if r.AccountID == 0 {
		return errors.ErrInvalidAccountID
	}
	// Amount is expected to be positive only, based on the operation type
	// it will be recorded accordingly
	if r.Amount <= 0 {
		return errors.ErrInvalidAmount
	}
	if !r.OperationTypeID.IsValid() {
		return errors.ErrInvalidOperationType
	}

	return nil
}

type CreateTransactionResponse struct {
	TransactionID int64                   `json:"transaction_id"`
	Status        types.TransactionStatus `json:"status"`
	StatusMessage string                  `json:"status_message,omitempty"`
}
