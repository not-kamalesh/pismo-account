package dto

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/not-kamalesh/pismo-account/errors"
)

type CreateAccountRequest struct {
	MsgID          string `json:"msg_id"`
	DocumentNumber string `json:"document_number"`
	Currency       string `json:"currency"`
}

func NewCreateAccountRequest() *CreateAccountRequest {
	return &CreateAccountRequest{}
}

func (r *CreateAccountRequest) Parse(httpReq *http.Request) error {
	decoder := json.NewDecoder(httpReq.Body)
	if err := decoder.Decode(r); err != nil {
		slog.Error("error occured to decode", "err", err)
		return err
	}

	return nil
}

func (r *CreateAccountRequest) Validate() error {
	if r.MsgID == "" {
		return errors.ErrInvalidMsgID
	}
	if r.DocumentNumber == "" {
		return errors.ErrInvalidDocumentID
	}
	if len(r.Currency) != 3 {
		return errors.ErrInvalidCurrency
	}
	return nil
}

type CreateAccountResponse struct {
	AccountID int64 `json:"account_id"`
}

type GetAccountRequest struct {
	MsgID     string `json:"msg_id"`
	AccountID int64  `json:"account_id"`
}

func NewGetAccountRequest() *GetAccountRequest {
	return &GetAccountRequest{}
}

func (r *GetAccountRequest) Parse(httpReq *http.Request) error {

	query := httpReq.URL.Query()
	if v := query.Get("msg_id"); v != "" {
		r.MsgID = v
	}

	vars := mux.Vars(httpReq)
	if v := vars["account_id"]; v != "" {
		accountID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		r.AccountID = accountID
	}

	return nil
}

func (r *GetAccountRequest) Validate() error {
	if r.MsgID == "" {
		return errors.ErrInvalidMsgID
	}
	if r.AccountID == 0 {
		return errors.ErrInvalidAccountID
	}
	return nil
}

type GetAccountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
	Currency       string `json:"currency"`
	Status         string `json:"status"`
}
