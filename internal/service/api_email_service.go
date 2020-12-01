package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/commitdev/zero-notification-service/internal/server"
)

// EmailApiService is a service that implents the logic for the EmailApiServicer
// This service should implement the business logic for every endpoint for the EmailApi API.
// Include any external packages or services that will be required by this service.
type EmailApiService struct {
}

// NewEmailApiService creates a default api service
func NewEmailApiService() server.EmailApiServicer {
	return &EmailApiService{}
}

// SendEmail - Send an email
func (s *EmailApiService) SendEmail(ctx context.Context, emailMessage server.EmailMessage) (server.ImplResponse, error) {
	// TODO - update SendEmail with the required logic for this service method.
	// Add api_email_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, string{}) or use other options such as http.Ok ...
	//return Response(200, string{}), nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return server.Response(http.StatusNotImplemented, nil), errors.New("SendEmail method not implemented")
}
