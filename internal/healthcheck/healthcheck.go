package healthcheck

import (
	"sync/atomic"

	"github.com/not-kamalesh/pismo-account/dto"
)

type Handler struct {
	isShuttingDown atomic.Bool
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HealthCheck() *dto.HeathCheckResponse {
	return &dto.HeathCheckResponse{
		Status:  "OK",
		Message: "I am Healthy, you chill",
	}
}
