package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/commitdev/zero-notification-service/internal/config"
	"github.com/commitdev/zero-notification-service/internal/server"
)

// EmailApiService is a service that implents the logic for the EmailApiServicer
// This service should implement the business logic for every endpoint for the EmailApi API.
// Include any external packages or services that will be required by this service.
type SmsApiService struct {
	config *config.Config
}

// NewSmsApiService creates a default api service
func NewSmsApiService(c *config.Config) server.SmsApiServicer {
	return &SmsApiService{c}
}

// SendSMS - Send an email
func (s *SmsApiService) SendSMS(ctx context.Context, sendSMSRequest server.SendSmsRequest) (server.ImplResponse, error) {
	// Set initial variables
	accountSid := s.config.TwilioAccID
	authToken := s.config.TwilioAuthToken
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	// Build out the data for our message
	v := url.Values{}
	v.Set("To", sendSMSRequest.Recipient)
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
	fmt.Print(resp.Status)

	if err != nil {
		fmt.Printf("Error sending SMS: %v\n", resp)
		return server.Response(http.StatusInternalServerError, nil), fmt.Errorf("Unable to send SMS: %v", err)
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		fmt.Printf("Failure from Twilio when sending SMS: %v", resp)
		return server.Response(http.StatusBadRequest, server.SendSmsResponse{Message: "Error sending SMS"}), nil
	}
	return server.Response(http.StatusOK, server.SendSmsResponse{Message: "SMS Sent"}), nil
}
