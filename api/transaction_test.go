package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/not-kamalesh/pismo-account/dto"
	"github.com/not-kamalesh/pismo-account/errors"
	"github.com/not-kamalesh/pismo-account/internal/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAPIHandler_CreateTransaction(t *testing.T) {
	tests := []struct {
		name           string
		request        *dto.CreateTransactionRequest
		setUpMocks     func(mockTransactionHandler *transaction.MockTransactionHandler)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "when request is empty, then write validation error",
			request:        &dto.CreateTransactionRequest{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":"INVALID_ARGUMENT","message":"Invalid argument provided"}`,
		},
		{
			name: "when request is invalid, then write validation error",
			request: &dto.CreateTransactionRequest{
				MsgID: "test_msgID",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":"INVALID_ARGUMENT","message":"Invalid argument provided"}`,
		},
		{
			name: "when request is valid and handler returns error, then write error",
			request: &dto.CreateTransactionRequest{
				MsgID:           "test_msgID",
				ReferenceID:     "test_referenceID",
				AccountID:       1,
				Amount:          100,
				OperationTypeID: 1,
			},
			setUpMocks: func(mockTransactionHandler *transaction.MockTransactionHandler) {
				mockTransactionHandler.On("Create", mock.Anything, mock.Anything).Return(nil, errors.ErrAlreadyExists).Once()
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"code":"ALREADY_EXISTS","message":"Resource already exists"}`,
		},
		{
			name: "when request is valid and handler returns response, then write response",
			request: &dto.CreateTransactionRequest{
				MsgID:           "test_msgID",
				ReferenceID:     "test_referenceID",
				AccountID:       1,
				Amount:          100,
				OperationTypeID: 1,
			},
			setUpMocks: func(mockTransactionHandler *transaction.MockTransactionHandler) {
				mockTransactionHandler.On("Create", mock.Anything, mock.Anything).Return(&dto.CreateTransactionResponse{TransactionID: 1}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"transaction_id":1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare httpRequest and response writer
			reqJson, _ := json.Marshal(tt.request)
			httpReq := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(reqJson))
			respWriter := httptest.NewRecorder()

			// setup mocks, create handler and execute the handler
			mockTransactionHandler := new(transaction.MockTransactionHandler)
			if tt.setUpMocks != nil {
				tt.setUpMocks(mockTransactionHandler)
			}
			api := NewAPIHandler(nil, nil, mockTransactionHandler)
			api.CreateTransaction(respWriter, httpReq)

			// assert the expectations
			res := respWriter.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.JSONEq(t, tt.expectedBody, string(body))
			mockTransactionHandler.AssertExpectations(t)
		})
	}
}
