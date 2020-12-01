package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/commitdev/zero-notification-service/internal/config"
	"github.com/commitdev/zero-notification-service/internal/server"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailApiService is a service that implents the logic for the EmailApiServicer
// This service should implement the business logic for every endpoint for the EmailApi API.
// Include any external packages or services that will be required by this service.
type EmailApiService struct {
	config *config.Config
}

// NewEmailApiService creates a default api service
func NewEmailApiService(c *config.Config) server.EmailApiServicer {
	return &EmailApiService{c}
}

// SendEmail - Send an email
func (s *EmailApiService) SendEmail(ctx context.Context, emailMessage server.EmailMessage) (server.ImplResponse, error) {
	from := mail.NewEmail(emailMessage.To.Name, emailMessage.To.Address)
	to := mail.NewEmail(emailMessage.From.Name, emailMessage.From.Address)
	message := mail.NewSingleEmail(from, emailMessage.Subject, to, emailMessage.Body, emailMessage.RichBody)
	client := sendgrid.NewSendClient(s.config.SendgridAPIKey)
	response, err := client.Send(message)
	fmt.Printf("%v", response)
	if err != nil {
		return server.Response(http.StatusInternalServerError, nil), fmt.Errorf("Unable to send email: %v", err)
	}

	return server.Response(http.StatusOK, server.SendMailResponse{}), nil
}
