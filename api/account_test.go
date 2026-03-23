package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/not-kamalesh/pismo-account/dto"
	"github.com/not-kamalesh/pismo-account/errors"
	"github.com/not-kamalesh/pismo-account/internal/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAPIHandler_CreateAccount(t *testing.T) {
	tests := []struct {
		name           string
		request        *dto.CreateAccountRequest
		setUpMocks     func(mock *account.MockAccountHandler)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "when request is empty, then write validation error",
			request:        &dto.CreateAccountRequest{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":"INVALID_ARGUMENT","message":"Invalid msg_id"}`,
		},
		{
			name: "when request is invalid, then write validation error",
			request: &dto.CreateAccountRequest{
				MsgID: "test_msgID",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":"INVALID_ARGUMENT","message":"Invalid document_id"}`,
		},
		{
			name: "when request is valid and handler returns error, then write error",
			request: &dto.CreateAccountRequest{
				MsgID:          "test_msgID",
				DocumentNumber: "1234",
				Currency:       "INR",
			},
			setUpMocks: func(mockAccount *account.MockAccountHandler) {
				mockAccount.On("Create", mock.Anything, mock.Anything).Return(nil, errors.ErrInternal).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":"INTERNAL","message":"Internal server error"}`,
		},
		{
			name: "when request is valid and handler returns response, then write response",
			request: &dto.CreateAccountRequest{
				MsgID:          "test_msgID",
				DocumentNumber: "1234",
				Currency:       "INR",
			},
			setUpMocks: func(mockAccount *account.MockAccountHandler) {
				mockAccount.On("Create", mock.Anything, mock.Anything).Return(&dto.CreateAccountResponse{AccountID: 1}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"account_id":1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare httpRequest and response writer
			reqJson, _ := json.Marshal(tt.request)
			httpReq := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(reqJson))
			respWriter := httptest.NewRecorder()

			// setup mocks, create handler and execute the handler
			mockAccountHandler := new(account.MockAccountHandler)
			if tt.setUpMocks != nil {
				tt.setUpMocks(mockAccountHandler)
			}
			api := NewAPIHandler(nil, mockAccountHandler, nil)
			api.CreateAccount(respWriter, httpReq)

			// assert the expectations
			res := respWriter.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.JSONEq(t, tt.expectedBody, string(body))

			mockAccountHandler.AssertExpectations(t)
		})
	}
}

func TestAPIHandler_GetAccount(t *testing.T) {
	tests := []struct {
		name           string
		request        *dto.GetAccountRequest
		setUpMocks     func(mock *account.MockAccountHandler)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "when request is empty, then write validation error",
			request:        &dto.GetAccountRequest{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":"INVALID_ARGUMENT","message":"Invalid msg_id"}`,
		},
		{
			name: "when request is invalid, then write validation error",
			request: &dto.GetAccountRequest{
				MsgID: "test_msgID",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":"INVALID_ARGUMENT","message":"Invalid account_id"}`,
		},
		{
			name: "when request is valid and handler returns error, then write error",
			request: &dto.GetAccountRequest{
				MsgID:     "test_msgID",
				AccountID: 1,
			},
			setUpMocks: func(mockAccount *account.MockAccountHandler) {
				mockAccount.On("Get", mock.Anything, mock.Anything).Return(nil, errors.ErrAlreadyExists).Once()
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"code":"ALREADY_EXISTS","message":"Resource already exists"}`,
		},
		{
			name: "when request is valid and handler returns response, then write response",
			request: &dto.GetAccountRequest{
				MsgID:     "test_msgID",
				AccountID: 1,
			},
			setUpMocks: func(mockAccount *account.MockAccountHandler) {
				mockAccount.On("Get", mock.Anything, mock.Anything).Return(&dto.GetAccountResponse{
					AccountID:      1,
					DocumentNumber: "1234",
					Currency:       "INR",
					Status:         "ACTIVE",
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"account_id":1, "document_number": "1234", "currency" : "INR", "status": "ACTIVE"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare httpRequest and response writer
			url := "/accounts/0"
			if tt.request.AccountID != 0 {
				url = fmt.Sprintf("/accounts/%d", tt.request.AccountID)
			}
			if tt.request.MsgID != "" {
				url = fmt.Sprintf("%s?msg_id=%s", url, tt.request.MsgID)
			}
			httpReq := httptest.NewRequest(http.MethodGet, url, nil)
			httpReq = mux.SetURLVars(httpReq, map[string]string{
				"account_id": fmt.Sprintf("%d", tt.request.AccountID),
			})
			respWriter := httptest.NewRecorder()

			// setup mocks, create handler and execute the handler
			mockAccountHandler := new(account.MockAccountHandler)
			if tt.setUpMocks != nil {
				tt.setUpMocks(mockAccountHandler)
			}
			api := NewAPIHandler(nil, mockAccountHandler, nil)
			api.GetAccount(respWriter, httpReq)

			// assert the expectations
			res := respWriter.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.JSONEq(t, tt.expectedBody, string(body))

			mockAccountHandler.AssertExpectations(t)
		})
	}
}
