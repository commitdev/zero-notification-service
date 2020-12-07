package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/commitdev/zero-notification-service/internal/config"
	"github.com/commitdev/zero-notification-service/internal/mail"
	"github.com/commitdev/zero-notification-service/internal/server"
	"github.com/sendgrid/sendgrid-go"
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
func (s *EmailApiService) SendEmail(ctx context.Context, sendMailRequest server.SendMailRequest) (server.ImplResponse, error) {
	client := sendgrid.NewSendClient(s.config.SendgridAPIKey)
	response, err := mail.SendIndividualMail(sendMailRequest.To, sendMailRequest.From, sendMailRequest.Message, client)

	if err != nil {
		fmt.Printf("Error sending mail: %v\n", response)
		return server.Response(http.StatusInternalServerError, nil), fmt.Errorf("Unable to send email: %v", err)
	}

	if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
		fmt.Printf("Failure from Sendgrid when sending mail: %v\n", response)
		return server.Response(http.StatusInternalServerError, nil), fmt.Errorf("Unable to send email: %v from mail provider: %v", response.StatusCode, response.Body)
	}

	return server.Response(http.StatusOK, server.SendMailResponse{response.Headers["X-Message-Id"][0]}), nil
}

// SendBulk - Send a batch of emails to many users with the same content. Note that it is possible for only a subset of these to fail.
func (s *EmailApiService) SendBulk(ctx context.Context, sendBulkMailRequest server.SendBulkMailRequest) (server.ImplResponse, error) {
	client := sendgrid.NewSendClient(s.config.SendgridAPIKey)

	responseChannel := make(chan mail.BulkSendAttempt)

	mail.SendBulkMail(sendBulkMailRequest.ToAddresses, sendBulkMailRequest.From, sendBulkMailRequest.Message, client, responseChannel)

	var successful []server.SendBulkMailResponseSuccessful
	var failed []server.SendBulkMailResponseFailed

	// Read all the responses from the channel. This will block if responses aren't ready and the channel is not yet closed
	for r := range responseChannel {
		if r.Error != nil {
			fmt.Printf("Error sending bulk mail: %v", r.Error)
			failed = append(failed, server.SendBulkMailResponseFailed{r.EmailAddress, fmt.Sprintf("Unable to send email: %v\n", r.Error)})
		} else if !(r.Response.StatusCode >= 200 && r.Response.StatusCode <= 299) {
			fmt.Printf("Failure from Sendgrid when sending bulk mail: %v", r.Response)
			failed = append(failed, server.SendBulkMailResponseFailed{r.EmailAddress, fmt.Sprintf("Unable to send email: %v from mail provider: %v\n", r.Response.StatusCode, r.Response.Body)})
		} else {
			successful = append(successful, server.SendBulkMailResponseSuccessful{r.EmailAddress, r.Response.Headers["X-Message-Id"][0]})
		}
	}
	responseCode := http.StatusOK
	if len(successful) == 0 {
		responseCode = http.StatusInternalServerError
	}
	return server.Response(responseCode, server.SendBulkMailResponse{successful, failed}), nil
}
