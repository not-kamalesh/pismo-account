package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/not-kamalesh/pismo-account/common/types"
	"github.com/not-kamalesh/pismo-account/dto"
	"github.com/not-kamalesh/pismo-account/internal/idempotencymgr"
	"github.com/not-kamalesh/pismo-account/internal/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ConcurrentCreateTransaction(t *testing.T) {
	type request struct {
		req            *dto.CreateTransactionRequest
		expectedStatus int
		expectedBody   string
	}
	tests := []struct {
		name           string
		request1       *request
		request2       *request
		concurrency    int
		isConflictCase bool
		setUpMocks     func(mockTransactionHandler *transaction.MockTransactionHandler)
	}{
		{
			name: "when same request tried multiple times, then it should produce the same result",
			request1: &request{
				req: &dto.CreateTransactionRequest{
					MsgID:           "test_msgid_1",
					ReferenceID:     "test_reference_id_1",
					AccountID:       1,
					Amount:          12.00,
					OperationTypeID: 1,
				},
				expectedStatus: 200,
				expectedBody:   `{"transaction_id":1, "status": "success"}`,
			},
			concurrency: 100,
			setUpMocks: func(mockTransactionHandler *transaction.MockTransactionHandler) {
				// should be called only once, other threads should wait
				// for the first thread to complete processing
				mockTransactionHandler.On("Create", mock.Anything, mock.Anything).
					Return(&dto.CreateTransactionResponse{
						TransactionID: 1,
						Status:        types.Success,
					}, nil).Once()

			},
		},
		{
			name: "when two different request tried multiple times, then it should produce the result based on the request",
			request1: &request{
				req: &dto.CreateTransactionRequest{
					MsgID:           "test_msgid_1",
					ReferenceID:     "test_reference_id_1",
					AccountID:       1,
					Amount:          12.00,
					OperationTypeID: 1,
				},
				expectedStatus: 200,
				expectedBody:   `{"transaction_id":1, "status": "success"}`,
			},
			request2: &request{
				req: &dto.CreateTransactionRequest{
					MsgID:           "test_msgid_2",
					ReferenceID:     "test_reference_id_2",
					AccountID:       1,
					Amount:          15.00,
					OperationTypeID: 1,
				},
				expectedStatus: 200,
				expectedBody:   `{"transaction_id":2, "status": "success"}`,
			},
			concurrency: 100,
			setUpMocks: func(mockTransactionHandler *transaction.MockTransactionHandler) {
				mockTransactionHandler.
					On("Create", mock.Anything, mock.MatchedBy(func(req *dto.CreateTransactionRequest) bool {
						return req.ReferenceID == "test_reference_id_1"
					})).
					Return(&dto.CreateTransactionResponse{
						TransactionID: 1,
						Status:        types.Success,
					}, nil).Once()
				mockTransactionHandler.
					On("Create", mock.Anything, mock.MatchedBy(func(req *dto.CreateTransactionRequest) bool {
						return req.ReferenceID == "test_reference_id_2"
					})).
					Return(&dto.CreateTransactionResponse{
						TransactionID: 2,
						Status:        types.Success,
					}, nil).Once()

			},
		},
		{
			name: "when same request(idempotency_key same) with different body tried multiple times, then it should produce return conflict for the second received request consistantly",
			request1: &request{
				req: &dto.CreateTransactionRequest{
					MsgID:           "test_msgid_1",
					ReferenceID:     "test_reference_id_1",
					AccountID:       1,
					Amount:          12.00,
					OperationTypeID: 1,
				},
				expectedStatus: 200,
				expectedBody:   `{"transaction_id":1, "status": "success"}`,
			},
			request2: &request{
				req: &dto.CreateTransactionRequest{
					MsgID:           "test_msgid_2",
					ReferenceID:     "test_reference_id_1",
					AccountID:       1,
					Amount:          15.00,
					OperationTypeID: 1,
				},
				expectedStatus: 409,
				expectedBody:   `{"code":"CONFLICT", "message": "Transaction Processed with different payload"}`,
			},
			concurrency:    100,
			isConflictCase: true,
			setUpMocks: func(mockTransactionHandler *transaction.MockTransactionHandler) {
				mockTransactionHandler.
					On("Create", mock.Anything, mock.MatchedBy(func(req *dto.CreateTransactionRequest) bool {
						return req.ReferenceID == "test_reference_id_1"
					})).
					Return(&dto.CreateTransactionResponse{
						TransactionID: 1,
						Status:        types.Success,
					}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup mocks, create handler and execute the handler
			mockTransactionHandler := new(transaction.MockTransactionHandler)
			idempotencyMgr := idempotencymgr.NewInMemIdempotencyMgr()
			if tt.setUpMocks != nil {
				tt.setUpMocks(mockTransactionHandler)
			}
			api := NewAPIHandler(nil, nil, mockTransactionHandler, idempotencyMgr)
			executeOneRequest := func(request *request) {
				// prepare httpRequest and response writer
				reqJson, _ := json.Marshal(request.req)
				httpReq := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(reqJson))
				respWriter := httptest.NewRecorder()

				api.CreateTransaction(respWriter, httpReq)

				// assert the expectations
				res := respWriter.Result()
				defer res.Body.Close()
				body, _ := io.ReadAll(res.Body)

				assert.Equal(t, request.expectedStatus, res.StatusCode)
				assert.JSONEq(t, request.expectedBody, string(body))
			}

			// Concurrenct request execution

			if tt.isConflictCase {
				// execute the request 1 first, just to say that request type 1 was processed first(to assert deterministically)
				executeOneRequest(tt.request1)
			}
			var wg sync.WaitGroup
			for i := 0; i < tt.concurrency; i++ {
				if tt.request1 != nil {
					wg.Add(1)
					go func() {
						defer wg.Done()
						executeOneRequest(tt.request1)
					}()
				}
				if tt.request2 != nil {
					wg.Add(1)
					go func() {
						defer wg.Done()
						executeOneRequest(tt.request2)
					}()
				}
			}
			wg.Wait()
		})
	}
}
