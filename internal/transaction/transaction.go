package transaction

import (
	"context"

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

func (h *transactionHandler) Create(ctx context.Context, req *dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error) {

	// Load the transaction based on reference_id, if exists return txID
	txn, loadErr := h.transactionDAO.LoadByReferenceID(ctx, req.ReferenceID)
	if loadErr != nil && loadErr != gorm.ErrRecordNotFound {
		return nil, loadErr
	}
	if txn != nil {
		return &dto.CreateTransactionResponse{
			TransactionID: txn.ID,
		}, nil
	}

	// Load the account based on account ID
	account, loadErr := h.accountDAO.LoadByID(ctx, req.AccountID)
	if loadErr != nil {
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
		return nil, saveErr
	}

	return &dto.CreateTransactionResponse{
		TransactionID: newTxn.ID,
	}, nil
}
