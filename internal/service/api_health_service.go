package service

import (
	"context"
	"net/http"

	"github.com/commitdev/zero-notification-service/internal/server"
)

// HealthApiService is a service that implents the logic for the HealthApiServicer
// This service should implement the business logic for every endpoint for the HealthApi API.
// Include any external packages or services that will be required by this service.
type HealthApiService struct {
}

// NewHealthApiService creates a default api service
func NewHealthApiService() server.HealthApiServicer {
	return &server.HealthApiService{}
}

// ReadyCheck - Readiness check - the service is ready to handle work
func (s *HealthApiService) ReadyCheck(ctx context.Context) (server.ImplResponse, error) {
	return server.Response(http.StatusOK, ""), nil
}
