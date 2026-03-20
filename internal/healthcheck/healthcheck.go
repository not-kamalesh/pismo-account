package healthcheck

import (
	"sync/atomic"

	"github.com/not-kamalesh/pismo-account/dto"
)

//go:generate mockery --name=HealthCheckHandler --output=. --outpkg=healthcheck --filename=mock_healthcheck.go --structname=MockHealthCheck
type HealthCheckHandler interface {
	HealthCheck() *dto.HeathCheckResponse
}

type healthCheckHandler struct {
	isShuttingDown atomic.Bool
}

func NewHandler() HealthCheckHandler {
	return &healthCheckHandler{}
}

func (h *healthCheckHandler) HealthCheck() *dto.HeathCheckResponse {
	return &dto.HeathCheckResponse{
		Status:  "OK",
		Message: "I am Healthy, you chill",
	}
}
