package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/commitdev/zero-notification-service/internal/config"
	"github.com/commitdev/zero-notification-service/internal/mail"
	"github.com/commitdev/zero-notification-service/internal/server"
	"github.com/sendgrid/sendgrid-go"
	"go.uber.org/zap"
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

	// Check if there are valid recipients who are not in the restriction list
	if len(s.config.AllowEmailToDomains) > 0 {
		originalAddresses := sendMailRequest.ToAddresses
		sendMailRequest.ToAddresses = mail.RemoveInvalidRecipients(sendMailRequest.ToAddresses, s.config.AllowEmailToDomains)
		// If there are none, return with a 200 but don't send anything
		if len(sendMailRequest.ToAddresses) == 0 {
			zap.S().Infow("No valid Recipients for send", zap.Any("original_addresses", originalAddresses))
			return server.Response(http.StatusOK, server.SendMailResponse{TrackingId: "No valid recipients"}), nil
		}

		sendMailRequest.CcAddresses = mail.RemoveInvalidRecipients(sendMailRequest.CcAddresses, s.config.AllowEmailToDomains)
		sendMailRequest.BccAddresses = mail.RemoveInvalidRecipients(sendMailRequest.BccAddresses, s.config.AllowEmailToDomains)
	}

	client := sendgrid.NewSendClient(s.config.SendgridAPIKey)
	response, err := mail.SendIndividualMail(sendMailRequest.ToAddresses, sendMailRequest.FromAddress, sendMailRequest.CcAddresses, sendMailRequest.BccAddresses, sendMailRequest.Message, client)

	if err != nil {
		zap.S().Errorf("Error sending mail: %v", response)

		return server.Response(http.StatusInternalServerError, nil), fmt.Errorf("Unable to send email: %v", err)
	}

	if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
		zap.S().Errorf("Failure from Sendgrid when sending mail: %v", response)
		return server.Response(http.StatusInternalServerError, nil), fmt.Errorf("Unable to send email: %v from mail provider: %v", response.StatusCode, response.Body)
	}

	return server.Response(http.StatusOK, server.SendMailResponse{TrackingId: response.Headers["X-Message-Id"][0]}), nil
}

// SendBulk - Send a batch of emails to many users with the same content. Note that it is possible for only a subset of these to fail.
func (s *EmailApiService) SendBulk(ctx context.Context, sendBulkMailRequest server.SendBulkMailRequest) (server.ImplResponse, error) {
	// Check if there are valid recipients who are not in the restriction list
	if len(s.config.AllowEmailToDomains) > 0 {
		originalAddresses := sendBulkMailRequest.ToAddresses
		sendBulkMailRequest.ToAddresses = mail.RemoveInvalidRecipients(sendBulkMailRequest.ToAddresses, s.config.AllowEmailToDomains)

		// If there are none, return with a 200 but don't send anything
		if len(sendBulkMailRequest.ToAddresses) == 0 {
			zap.S().Infow("No valid Recipients for bulk send", zap.Any("original_addresses", originalAddresses))
			return server.Response(http.StatusOK, server.SendMailResponse{TrackingId: "No valid recipients"}), nil
		}

		sendBulkMailRequest.CcAddresses = mail.RemoveInvalidRecipients(sendBulkMailRequest.CcAddresses, s.config.AllowEmailToDomains)
		sendBulkMailRequest.BccAddresses = mail.RemoveInvalidRecipients(sendBulkMailRequest.BccAddresses, s.config.AllowEmailToDomains)
	}

	client := sendgrid.NewSendClient(s.config.SendgridAPIKey)

	responseChannel := make(chan mail.BulkSendAttempt)

	mail.SendBulkMail(sendBulkMailRequest.ToAddresses, sendBulkMailRequest.FromAddress, sendBulkMailRequest.CcAddresses, sendBulkMailRequest.BccAddresses, sendBulkMailRequest.Message, client, responseChannel)

	var successful []server.SendBulkMailResponseSuccessful
	var failed []server.SendBulkMailResponseFailed

	// Read all the responses from the channel. This will block if responses aren't ready and the channel is not yet closed
	for r := range responseChannel {
		if r.Error != nil {
			zap.S().Errorf("Error sending bulk mail: %v", r.Error)
			failed = append(failed, server.SendBulkMailResponseFailed{EmailAddress: r.EmailAddress, Error: fmt.Sprintf("Unable to send email: %v\n", r.Error)})
		} else if !(r.Response.StatusCode >= 200 && r.Response.StatusCode <= 299) {
			zap.S().Errorf("Failure from Sendgrid when sending bulk mail: %v", r.Response)
			failed = append(failed, server.SendBulkMailResponseFailed{EmailAddress: r.EmailAddress, Error: fmt.Sprintf("Unable to send email: %v from mail provider: %v\n", r.Response.StatusCode, r.Response.Body)})
		} else {
			successful = append(successful, server.SendBulkMailResponseSuccessful{EmailAddress: r.EmailAddress, TrackingId: r.Response.Headers["X-Message-Id"][0]})
		}
	}
	responseCode := http.StatusOK
	if len(successful) == 0 {
		responseCode = http.StatusBadRequest
	}
	return server.Response(responseCode, server.SendBulkMailResponse{Successful: successful, Failed: failed}), nil
}
