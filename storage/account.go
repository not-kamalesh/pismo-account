package storage

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID         int64 `gorm:"primaryKey"`
	DocumentID string
	Currency   string
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (a *Account) GetTableName() string {
	return "accounts"
}

//go:generate mockery --name=IAccountDao --inpackage --filename=mock_account_dao.go --structname=MockIAccountDAO

// IAccountDao : DAO interface for other modules to use
type IAccountDao interface {
	LoadByID(ctx context.Context, accountID int64) (*Account, error)
	LoadByDocumentID(ctx context.Context, documentID string) (*Account, error)
	Save(ctx context.Context, account *Account) error
	UpdateStatus(ctx context.Context, accountID int64, status string) error
}

type AccountDao struct {
	db *gorm.DB
}

func NewAccountDao(db *gorm.DB) IAccountDao {
	return &AccountDao{
		db: db,
	}
}

func (a *AccountDao) LoadByID(ctx context.Context, accountID int64) (*Account, error) {
	var account Account

	err := a.db.WithContext(ctx).Where("id = ?", accountID).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountDao) LoadByDocumentID(ctx context.Context, documentID string) (*Account, error) {
	var account Account

	err := a.db.WithContext(ctx).Where("document_id = ?", documentID).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountDao) Save(ctx context.Context, account *Account) error {

	err := a.db.WithContext(ctx).Save(account).Error
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountDao) UpdateStatus(ctx context.Context, accountID int64, status string) error {
	err := a.db.WithContext(ctx).Model(&Account{}).Where("id = ?", accountID).Update("status", status).Error
	if err != nil {
		return err
	}

	return nil
}
