package dto

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/not-kamalesh/pismo-account/common/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransactionRequest_Parse(t *testing.T) {
	tests := []struct {
		name    string
		body    interface{}
		want    *CreateTransactionRequest
		wantErr bool
	}{
		{
			name: "valid body",
			body: map[string]interface{}{
				"msg_id":            "123",
				"reference_id":      "ref-123",
				"account_id":        42,
				"operation_type_id": 1,
				"amount":            10.5,
			},
			want: &CreateTransactionRequest{
				MsgID:           "123",
				ReferenceID:     "ref-123",
				AccountID:       42,
				OperationTypeID: 1,
				Amount:          10.5,
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			body:    "{INVALID",
			want:    &CreateTransactionRequest{},
			wantErr: true,
		},
		{
			name:    "empty body",
			body:    nil,
			want:    &CreateTransactionRequest{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			case map[string]interface{}:
				bodyBytes, _ = json.Marshal(v)
			case nil:
				bodyBytes = nil
			}

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyBytes))
			assert.NoError(t, err)

			var r CreateTransactionRequest
			err = r.Parse(req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, &r)
			}
		})
	}
}

func TestCreateTransactionRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTransactionRequest
		wantErr bool
	}{
		{
			name: "valid",
			req: CreateTransactionRequest{
				MsgID:           "msg1",
				ReferenceID:     "ref1",
				AccountID:       5,
				OperationTypeID: types.OperationType(1),
				Amount:          20.0,
			},
			wantErr: false,
		},
		{
			name: "missing msg_id",
			req: CreateTransactionRequest{
				ReferenceID:     "ref2",
				AccountID:       5,
				OperationTypeID: types.OperationType(1),
				Amount:          20.0,
			},
			wantErr: true,
		},
		{
			name: "missing reference_id",
			req: CreateTransactionRequest{
				MsgID:           "msg2",
				AccountID:       5,
				OperationTypeID: types.OperationType(1),
				Amount:          20.0,
			},
			wantErr: true,
		},
		{
			name: "zero account_id",
			req: CreateTransactionRequest{
				MsgID:           "msg3",
				ReferenceID:     "ref3",
				OperationTypeID: types.OperationType(1),
				Amount:          20.0,
			},
			wantErr: true,
		},
		{
			name: "amount not positive",
			req: CreateTransactionRequest{
				MsgID:           "msg4",
				ReferenceID:     "ref4",
				AccountID:       5,
				OperationTypeID: types.OperationType(1),
				Amount:          0,
			},
			wantErr: true,
		},
		{
			name: "invalid operation_type",
			req: CreateTransactionRequest{
				MsgID:           "msg6",
				ReferenceID:     "ref6",
				AccountID:       5,
				OperationTypeID: types.OperationType(100),
				Amount:          10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
