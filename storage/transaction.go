package storage

import (
	"context"
	"time"

	"github.com/not-kamalesh/pismo-account/common/types"
	"gorm.io/gorm"
)

type Transaction struct {
	ID              int64  `gorm:"primaryKey"`
	ReferenceID     string `gorm:"type:varchar(64);size:64;uniqueIndex"`
	AccountID       int64
	OperationTypeID types.OperationType
	EntryType       types.EntryType
	Amount          int64
	Currency        string
	CreatedAt       time.Time
}

func (a *Transaction) GetTableName() string {
	return "transactions"
}

//go:generate mockery --name=ITransactionDao --inpackage --filename=mock_transaction_dao.go --structname=MockITransactionDAO

// ITransactionDao : DAO interface for other modules to use
type ITransactionDao interface {
	LoadByID(ctx context.Context, txnID int64) (*Transaction, error)
	LoadByReferenceID(ctx context.Context, referenceID string) (*Transaction, error)
	Save(ctx context.Context, txn *Transaction) error
}

type TransactionDao struct {
	db *gorm.DB
}

func NewTransactionDao(db *gorm.DB) ITransactionDao {
	return &TransactionDao{
		db: db,
	}
}

func (a *TransactionDao) LoadByID(ctx context.Context, txnID int64) (*Transaction, error) {
	var txn Transaction

	err := a.db.WithContext(ctx).Where("id = ?", txnID).First(&txn).Error
	if err != nil {
		return nil, err
	}

	return &txn, nil
}

func (a *TransactionDao) LoadByReferenceID(ctx context.Context, referenceID string) (*Transaction, error) {
	var txn Transaction

	err := a.db.WithContext(ctx).Where("reference_id = ?", referenceID).First(&txn).Error
	if err != nil {
		return nil, err
	}

	return &txn, nil
}

func (a *TransactionDao) Save(ctx context.Context, txn *Transaction) error {

	err := a.db.WithContext(ctx).Save(txn).Error
	if err != nil {
		return err
	}

	return nil
}
