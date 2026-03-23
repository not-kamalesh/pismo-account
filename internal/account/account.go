package account

import (
	"context"
	"log/slog"

	"github.com/not-kamalesh/pismo-account/dto"
	"github.com/not-kamalesh/pismo-account/storage"
	"gorm.io/gorm"
)

//go:generate mockery --name=AccountHandler --output=. --outpkg=account --filename=mock_account.go --structname=MockAccountHandler
type AccountHandler interface {
	Create(ctx context.Context, req *dto.CreateAccountRequest) (*dto.CreateAccountResponse, error)
	Get(ctx context.Context, req *dto.GetAccountRequest) (*dto.GetAccountResponse, error)
}

type accountHandler struct {
	dao storage.IAccountDao
}

func NewHandler(accountDAO storage.IAccountDao) AccountHandler {
	return &accountHandler{
		dao: accountDAO,
	}
}

func (h *accountHandler) Create(ctx context.Context, req *dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {

	// Check if an account already exists with the given document number
	account, loadErr := h.dao.LoadByDocumentID(ctx, req.DocumentNumber)
	if loadErr != nil && loadErr != gorm.ErrRecordNotFound {
		slog.Warn("error occured on loading account by documentID", "error", loadErr)
		return nil, loadErr
	}
	if account != nil {
		return &dto.CreateAccountResponse{
			AccountID: account.ID,
		}, nil
	}

	// Save the account record
	newAccount := &storage.Account{
		DocumentID: req.DocumentNumber,
		Currency:   req.Currency,
		Status:     "ACTIVE",
	}
	saveErr := h.dao.Save(ctx, newAccount)
	if saveErr != nil {
		slog.Warn("error occured on save account", "error", saveErr)
		return nil, saveErr
	}

	return &dto.CreateAccountResponse{
		AccountID: newAccount.ID,
	}, nil
}

func (h *accountHandler) Get(ctx context.Context, req *dto.GetAccountRequest) (*dto.GetAccountResponse, error) {

	account, loadErr := h.dao.LoadByID(ctx, req.AccountID)
	if loadErr != nil {
		slog.Warn("error occured on loading account by accountID", "error", loadErr)
		return nil, loadErr
	}

	return &dto.GetAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentID,
		Currency:       account.Currency,
		Status:         account.Status,
	}, nil
}
