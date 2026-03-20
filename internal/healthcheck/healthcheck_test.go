package healthcheck

import (
	"testing"

	"github.com/not-kamalesh/pismo-account/dto"
	"github.com/stretchr/testify/assert"
)

func TestHandlerHealthCheck(t *testing.T) {
	tests := []struct {
		name         string
		expectedResp *dto.HeathCheckResponse
	}{
		{
			name: "health check success",
			expectedResp: &dto.HeathCheckResponse{
				Status:  "OK",
				Message: "I am Healthy, you chill",
			},
		},
	}

	for _, scenario := range tests {
		t.Run(scenario.name, func(t *testing.T) {
			healthCheckHandler := NewHandler()
			resp := healthCheckHandler.HealthCheck()
			assert.Equal(t, scenario.expectedResp, resp)
		})
	}
}
