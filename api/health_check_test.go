package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/not-kamalesh/pismo-account/dto"
	"github.com/not-kamalesh/pismo-account/internal/healthcheck"
	"github.com/stretchr/testify/assert"
)

func TestAPIHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		setUpMocks     func(hcMock *healthcheck.MockHealthCheck)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "health check handler returns success response",
			setUpMocks: func(hcMock *healthcheck.MockHealthCheck) {
				hcMock.On("HealthCheck").Return(&dto.HeathCheckResponse{
					Status:  "OK",
					Message: "mocked",
				}).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"OK","message":"mocked"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// prepare httpRequest and response writer
			httpReq := httptest.NewRequest(http.MethodGet, "/health", nil)
			respWriter := httptest.NewRecorder()

			// setup mocks, create handler and execute the handler
			mockHC := new(healthcheck.MockHealthCheck)
			if tt.setUpMocks != nil {
				tt.setUpMocks(mockHC)
			}
			api := NewAPIHandler(mockHC, nil, nil)
			api.HealthCheck(respWriter, httpReq)

			// assert the expectations
			res := respWriter.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.JSONEq(t, tt.expectedBody, string(body))

			mockHC.AssertExpectations(t)
		})
	}
}
