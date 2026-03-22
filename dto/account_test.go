package dto

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/not-kamalesh/pismo-account/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccountRequest_Parse(t *testing.T) {
	tests := []struct {
		name    string
		body    []byte
		want    *CreateAccountRequest
		wantErr bool
	}{
		{
			name:    "when body is invalid JSON, then return error",
			body:    []byte(`{`),
			wantErr: true,
		},
		{
			name: "when body is valid JSON, then parse succeeds",
			body: []byte(`{"msg_id":"test_msgID","document_number":"1234","currency":"INR"}`),
			want: &CreateAccountRequest{
				MsgID:          "test_msgID",
				DocumentNumber: "1234",
				Currency:       "INR",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpReq := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(tt.body))
			got := NewCreateAccountRequest()
			err := got.Parse(httpReq)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateAccountRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateAccountRequest
		wantErr error
	}{
		{
			name:    "when request is empty, then return validation error",
			req:     &CreateAccountRequest{},
			wantErr: errors.ErrInvalidArgument,
		},
		{
			name: "when request is invalid, then return validation error",
			req: &CreateAccountRequest{
				MsgID: "test_msgID",
			},
			wantErr: errors.ErrInvalidArgument,
		},
		{
			name: "when currency length is not three, then return validation error",
			req: &CreateAccountRequest{
				MsgID:          "test_msgID",
				DocumentNumber: "1234",
				Currency:       "US",
			},
			wantErr: errors.ErrInvalidArgument,
		},
		{
			name: "when request is valid, then validate succeeds",
			req: &CreateAccountRequest{
				MsgID:          "test_msgID",
				DocumentNumber: "1234",
				Currency:       "INR",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestGetAccountRequest_Parse(t *testing.T) {
	tests := []struct {
		name         string
		accountIDVar string
		msgIDQuery   string
		want         *GetAccountRequest
		wantErr      bool
	}{
		{
			name:         "when account_id is not a number, then return error",
			accountIDVar: "something",
			msgIDQuery:   "test_msgID",
			wantErr:      true,
		},
		{
			name:         "when URL has msg_id and account_id, then parse succeeds",
			accountIDVar: "1",
			msgIDQuery:   "test_msgID",
			want: &GetAccountRequest{
				MsgID:     "test_msgID",
				AccountID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/accounts/%s", tt.accountIDVar)
			if tt.msgIDQuery != "" {
				url = fmt.Sprintf("%s?msg_id=%s", url, tt.msgIDQuery)
			}
			httpReq := httptest.NewRequest(http.MethodGet, url, nil)
			httpReq = mux.SetURLVars(httpReq, map[string]string{
				"account_id": tt.accountIDVar,
			})

			got := NewGetAccountRequest()
			err := got.Parse(httpReq)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetAccountRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *GetAccountRequest
		wantErr error
	}{
		{
			name:    "when request is empty, then return validation error",
			req:     &GetAccountRequest{},
			wantErr: errors.ErrInvalidArgument,
		},
		{
			name: "when request is invalid, then return validation error",
			req: &GetAccountRequest{
				MsgID: "test_msgID",
			},
			wantErr: errors.ErrInvalidArgument,
		},
		{
			name: "when request is valid, then validate succeeds",
			req: &GetAccountRequest{
				MsgID:     "test_msgID",
				AccountID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
