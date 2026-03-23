package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/not-kamalesh/pismo-account/dto"
	"github.com/not-kamalesh/pismo-account/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestTransactionHandler_Create(t *testing.T) {

	tests := []struct {
		name         string
		req          *dto.CreateTransactionRequest
		setupMocks   func(accountDao *storage.MockIAccountDAO, transactionDao *storage.MockITransactionDAO)
		expectedResp *dto.CreateTransactionResponse
		expectedErr  error
	}{
		{
			name: "when the transaction with reference_id already exists, then return the txnID",
			req: &dto.CreateTransactionRequest{
				ReferenceID:     "ref1",
				AccountID:       1,
				Amount:          200,
				OperationTypeID: 1,
			},
			setupMocks: func(accountDao *storage.MockIAccountDAO, transactionDao *storage.MockITransactionDAO) {
				transactionDao.On("LoadByReferenceID", mock.Anything, mock.Anything).
					Return(&storage.Transaction{ID: 123}, nil).Once()
			},
			expectedResp: &dto.CreateTransactionResponse{TransactionID: 123},
			expectedErr:  nil,
		},
		{
			name: "when DAO returns error(aprat from ErrRecordNotFound), then return error",
			req: &dto.CreateTransactionRequest{
				ReferenceID:     "ref2",
				AccountID:       42,
				Amount:          50,
				OperationTypeID: 1,
			},
			setupMocks: func(accountDao *storage.MockIAccountDAO, transactionDao *storage.MockITransactionDAO) {
				transactionDao.On("LoadByReferenceID", mock.Anything, mock.Anything).
					Return(nil, errors.New("db error")).Once()
			},
			expectedResp: nil,
			expectedErr:  errors.New("db error"),
		},
		{
			name: "when the account id not exists, then return error",
			req: &dto.CreateTransactionRequest{
				ReferenceID:     "ref3",
				AccountID:       999,
				Amount:          34.12,
				OperationTypeID: 2,
			},
			setupMocks: func(accountDao *storage.MockIAccountDAO, transactionDao *storage.MockITransactionDAO) {
				transactionDao.On("LoadByReferenceID", mock.Anything, mock.Anything).
					Return(nil, gorm.ErrRecordNotFound).Once()
				accountDao.On("LoadByID", mock.Anything, mock.Anything).
					Return(nil, errors.New("not found")).Once()
			},
			expectedResp: nil,
			expectedErr:  errors.New("not found"),
		},
		{
			name: "when the transaction.Save failes with a error, return error",
			req: &dto.CreateTransactionRequest{
				ReferenceID:     "ref4",
				AccountID:       11,
				Amount:          10,
				OperationTypeID: 1,
			},
			setupMocks: func(accountDao *storage.MockIAccountDAO, transactionDao *storage.MockITransactionDAO) {
				transactionDao.On("LoadByReferenceID", mock.Anything, mock.Anything).
					Return(nil, gorm.ErrRecordNotFound).Once()
				accountDao.On("LoadByID", mock.Anything, mock.Anything).
					Return(&storage.Account{
						ID:       11,
						Currency: "USD",
					}, nil).Once()
				transactionDao.On("Save", mock.Anything, mock.Anything).
					Return(errors.New("db save error")).Once()
			},
			expectedResp: nil,
			expectedErr:  errors.New("db save error"),
		},
		{
			name: "when its a happy path, return the txnID in response",
			req: &dto.CreateTransactionRequest{
				ReferenceID:     "new_ref",
				AccountID:       99,
				Amount:          109.42,
				OperationTypeID: 1,
			},
			setupMocks: func(accountDao *storage.MockIAccountDAO, transactionDao *storage.MockITransactionDAO) {
				transactionDao.On("LoadByReferenceID", mock.Anything, mock.Anything).
					Return(nil, gorm.ErrRecordNotFound).Once()
				accountDao.On("LoadByID", mock.Anything, int64(99)).
					Return(&storage.Account{
						ID:       1,
						Currency: "INR",
					}, nil).Once()
				transactionDao.On("Save", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					tx := args.Get(1).(*storage.Transaction)
					tx.ID = 1
				}).Return(nil).Once()
			},
			expectedResp: &dto.CreateTransactionResponse{TransactionID: 1},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountDao := new(storage.MockIAccountDAO)
			transactionDao := new(storage.MockITransactionDAO)

			if tt.setupMocks != nil {
				tt.setupMocks(accountDao, transactionDao)
			}

			h := NewHandler(accountDao, transactionDao)
			resp, err := h.Create(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				}
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}

			accountDao.AssertExpectations(t)
			transactionDao.AssertExpectations(t)
		})
	}
}
