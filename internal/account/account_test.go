package account

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"github.com/not-kamalesh/pismo-account/dto"
	"github.com/not-kamalesh/pismo-account/errors"
	"github.com/not-kamalesh/pismo-account/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountHandler_Create(t *testing.T) {
	ctx := context.Background()

	validReq := func() *dto.CreateAccountRequest {
		return &dto.CreateAccountRequest{
			MsgID:          "msg-1",
			DocumentNumber: "DOC-100",
			Currency:       "USD",
		}
	}

	tests := []struct {
		name    string
		req     *dto.CreateAccountRequest
		setup   func(m *storage.MockIAccountDAO)
		want    *dto.CreateAccountResponse
		wantErr error
	}{
		{
			name: "when document already exists, then ErrAlreadyExists",
			req:  validReq(),
			setup: func(m *storage.MockIAccountDAO) {
				m.On("LoadByDocumentID", mock.Anything, mock.Anything).Return(&storage.Account{
					ID:         1,
					DocumentID: "DOC-100",
					Currency:   "USD",
					Status:     "ACTIVE",
				}, nil)
			},
			want: &dto.CreateAccountResponse{AccountID: 1},
		},
		{
			name: "when LoadByDocumentID fails with non-not-found error, then propagate error",
			req:  validReq(),
			setup: func(m *storage.MockIAccountDAO) {
				m.On("LoadByDocumentID", mock.Anything, mock.Anything).Return(nil, errors.ErrInternal)
			},
			wantErr: errors.ErrInternal,
		},
		{
			name: "when Save fails, then propagate error",
			req:  validReq(),
			setup: func(m *storage.MockIAccountDAO) {
				m.On("LoadByDocumentID", mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
				m.On("Save", mock.Anything, mock.Anything).Return(errors.ErrInternal)
			},
			wantErr: errors.ErrInternal,
		},
		{
			name: "when document is new, then save and return account id",
			req:  validReq(),
			setup: func(m *storage.MockIAccountDAO) {
				m.On("LoadByDocumentID", mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
				m.On("Save", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					args.Get(1).(*storage.Account).ID = 1
				}).Return(nil)
			},
			want: &dto.CreateAccountResponse{AccountID: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDAO := storage.NewMockIAccountDAO(t)
			tt.setup(mockDAO)

			h := NewHandler(mockDAO)
			got, err := h.Create(ctx, tt.req)

			if tt.wantErr != nil {
				assert.Nil(t, got)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAccountHandler_Get(t *testing.T) {
	ctx := context.Background()

	validReq := func() *dto.GetAccountRequest {
		return &dto.GetAccountRequest{
			MsgID:     "msg-1",
			AccountID: 99,
		}
	}

	tests := []struct {
		name    string
		req     *dto.GetAccountRequest
		setup   func(m *storage.MockIAccountDAO)
		want    *dto.GetAccountResponse
		wantErr error
	}{
		{
			name: "when LoadByID fails, then propagate error",
			req:  validReq(),
			setup: func(m *storage.MockIAccountDAO) {
				m.On("LoadByID", mock.Anything, mock.Anything).Return(nil, errors.ErrInternal)
			},
			wantErr: errors.ErrInternal,
		},
		{
			name: "when LoadByID returns record not found, then propagate error",
			req:  validReq(),
			setup: func(m *storage.MockIAccountDAO) {
				m.On("LoadByID", mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "when account exists, then return details",
			req:  validReq(),
			setup: func(m *storage.MockIAccountDAO) {
				m.On("LoadByID", mock.Anything, mock.Anything).Return(&storage.Account{
					ID:         99,
					DocumentID: "DOC-100",
					Currency:   "USD",
					Status:     "ACTIVE",
				}, nil)
			},
			want: &dto.GetAccountResponse{
				AccountID:      99,
				DocumentNumber: "DOC-100",
				Currency:       "USD",
				Status:         "ACTIVE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDAO := storage.NewMockIAccountDAO(t)
			tt.setup(mockDAO)

			h := NewHandler(mockDAO)
			got, err := h.Get(ctx, tt.req)

			if tt.wantErr != nil {
				assert.Nil(t, got)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
