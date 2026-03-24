package transaction

import (
	"context"
	"log/slog"

	"github.com/not-kamalesh/pismo-account/common/types"
	"github.com/not-kamalesh/pismo-account/dto"
	"github.com/not-kamalesh/pismo-account/storage"
	"github.com/not-kamalesh/pismo-account/utils/amount"
	"gorm.io/gorm"
)

//go:generate mockery --name=TransactionHandler --output=. --outpkg=transaction --filename=mock_transaction.go --structname=MockTransactionHandler
type TransactionHandler interface {
	Create(ctx context.Context, req *dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error)
}

type transactionHandler struct {
	accountDAO     storage.IAccountDao
	transactionDAO storage.ITransactionDao
}

func NewHandler(accountDAO storage.IAccountDao, transactionDao storage.ITransactionDao) TransactionHandler {
	return &transactionHandler{
		accountDAO:     accountDAO,
		transactionDAO: transactionDao,
	}
}

// Create - creates a Transaction based on the request
// Check if transaction already exists based on reference_id - second layer of idempotent behaviour
// if exists : return the existing transaction
// if notExists:
//   - Load account by account id
//   - Insert Transaction in DB and return the created transaction
func (h *transactionHandler) Create(ctx context.Context, req *dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error) {

	// Load the transaction based on reference_id, if exists return txID
	txn, loadErr := h.transactionDAO.LoadByReferenceID(ctx, req.ReferenceID)
	// any db error apart from ErrRecordNotFound, return error
	if loadErr != nil && loadErr != gorm.ErrRecordNotFound {
		slog.Warn("error occured on loading transaction by reference_id", "reference_id", req.ReferenceID, "error", loadErr)
		return nil, loadErr
	}
	// if transaction exists, return result
	if txn != nil {
		slog.Info("transaction already exists", "reference_id", req.ReferenceID)
		return &dto.CreateTransactionResponse{
			TransactionID: txn.ID,
			Status:        types.Success,
		}, nil
	}

	// When transaction not exists, create it

	// Load the account based on account ID
	account, loadErr := h.accountDAO.LoadByID(ctx, req.AccountID)
	if loadErr != nil {
		if loadErr == gorm.ErrRecordNotFound {
			return &dto.CreateTransactionResponse{
				Status:        types.Failed,
				StatusMessage: "Account Not Found",
			}, nil
		}
		slog.Warn("error occured on loading account by account_id", "account_id", req.AccountID, "error", loadErr)
		return nil, loadErr
	}

	amt := amount.NewAmountFromFloat(req.Amount, account.Currency)

	// Create a transaction record and save it
	newTxn := &storage.Transaction{
		ReferenceID:     req.ReferenceID,
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		EntryType:       req.OperationTypeID.GetEntryType(),
		Amount:          amt.ToMinorUnit(),
		Currency:        account.Currency,
	}

	saveErr := h.transactionDAO.Save(ctx, newTxn)
	if saveErr != nil {
		slog.Warn("error occured on saving new transaction", "reference_id", req.ReferenceID, "error", saveErr)
		return nil, saveErr
	}

	return &dto.CreateTransactionResponse{
		TransactionID: newTxn.ID,
		Status:        types.Success,
	}, nil
}
