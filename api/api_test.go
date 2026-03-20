package api

import (
	commonerrors "errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/not-kamalesh/pismo-account/dto"
	"github.com/not-kamalesh/pismo-account/errors"
	"github.com/not-kamalesh/pismo-account/internal/healthcheck"
	"github.com/stretchr/testify/assert"
)

func TestNewAPIHandler(t *testing.T) {
	healthCheckHandler := &healthcheck.MockHealthCheck{}
	apiHandler := NewAPIHandler(healthCheckHandler)
	assert.NotNil(t, apiHandler)
}

func TestAPIHandler_writeResponse(t *testing.T) {
	tests := []struct {
		name           string
		apiName        string
		apiResp        interface{}
		apiErr         error
		isNilResp      bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "when TestAPI returns a non-nil resp with nil error, when the isNilResp is false, then statuscode should be 200",
			apiName: "TestAPI",
			apiResp: &dto.HeathCheckResponse{
				Status:  "OK",
				Message: "test response",
			},
			apiErr:         nil,
			isNilResp:      false,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"OK","message":"test response"}`,
		},
		{
			name:           "when TestAPI returns a nil resp with nil error, when the isNilResp is true, then statuscode should be 204",
			apiName:        "TestAPI",
			apiResp:        nil,
			apiErr:         nil,
			isNilResp:      true,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "when TestAPI returns a nil resp with nil error, when the isNilResp is false, then statuscode should be 503",
			apiName:        "TestAPI",
			apiResp:        nil,
			apiErr:         nil,
			isNilResp:      false,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "when TestAPI returns a nil resp with non-nil error, when the isNilResp is false, then statuscode should be based on error",
			apiName:        "TestAPI",
			apiResp:        nil,
			apiErr:         errors.ErrInvalidArgument,
			isNilResp:      false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":"INVALID_ARGUMENT","message":"Invalid argument provided"}`,
		},
		{
			name:           "when TestAPI returns a nil resp with non-nil error, when the isNilResp is true, then statuscode should be based on error",
			apiName:        "TestAPI",
			apiResp:        nil,
			apiErr:         errors.ErrInvalidArgument,
			isNilResp:      true,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":"INVALID_ARGUMENT","message":"Invalid argument provided"}`,
		},
		{
			name:           "when TestAPI returns a nil resp with non-nil error(not pre-defined), when the isNilResp is true, then statuscode should be based on error",
			apiName:        "TestAPI",
			apiResp:        nil,
			apiErr:         commonerrors.New("some random error"),
			isNilResp:      true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":"INTERNAL","message":"some random error"}`,
		},
	}

	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {

			respWriter := httptest.NewRecorder()
			api := NewAPIHandler(nil)
			api.writeResponse(respWriter, scenario.apiName, scenario.apiResp, scenario.apiErr, scenario.isNilResp)

			// assert the expectations
			res := respWriter.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, scenario.expectedStatus, res.StatusCode)
			if scenario.expectedBody != "" {
				assert.JSONEq(t, scenario.expectedBody, string(body))
			}
		})
	}
}
