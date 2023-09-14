package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/commitdev/zero-notification-service/internal/config"
	"github.com/commitdev/zero-notification-service/internal/server"
	"go.uber.org/zap"
)

// EmailApiService is a service that implents the logic for the EmailApiServicer
// This service should implement the business logic for every endpoint for the EmailApi API.
// Include any external packages or services that will be required by this service.
type SmsApiService struct {
	config *config.Config
}

// NewSmsApiService creates a default api service
func NewSmsApiService(c *config.Config) server.SmsAPIServicer {
	return &SmsApiService{c}
}

// SendSMS - Send an email
func (s *SmsApiService) SendSMS(ctx context.Context, sendSMSRequest server.SendSmsRequest) (server.ImplResponse, error) {
	// Set initial variables
	accountSid := s.config.TwilioAccountID
	authToken := s.config.TwilioAuthToken
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	// Build out the data for our message
	v := url.Values{}
	v.Set("To", sendSMSRequest.RecipientPhoneNumber)
	v.Set("From", s.config.TwilioPhoneNumber)
	v.Set("Body", sendSMSRequest.Message)
	rb := *strings.NewReader(v.Encode())

	// Create Client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := client.Do(req)

	// Convert request.body response to string
	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	if err != nil {
		zap.S().Errorf("Error sending SMS: %v", err)
		return server.Response(http.StatusInternalServerError, nil), fmt.Errorf("unable to send SMS: %v", err)
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		zap.S().Errorf("Failure from Twilio when sending SMS: %v\n", bodyString)
		return server.Response(http.StatusBadRequest, server.SendSmsResponse{Message: fmt.Sprintf("Error sending SMS: %d\n, Message from SMS Provider: %v\n", resp.StatusCode, bodyString)}), nil
	}
	return server.Response(http.StatusOK, server.SendSmsResponse{Message: "SMS Sent"}), nil
}
