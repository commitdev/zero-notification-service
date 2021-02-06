package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/commitdev/zero-notification-service/internal/config"
	"github.com/commitdev/zero-notification-service/internal/notification/slack"
	"github.com/commitdev/zero-notification-service/internal/server"
	slack_lib "github.com/slack-go/slack"
)

// NotificationApiService is a service that implents the logic for the NotificationApiServicer
// This service should implement the business logic for every endpoint for the NotificationApi API.
// Include any external packages or services that will be required by this service.
type NotificationApiService struct {
	config *config.Config
}

// NewNotificationApiService creates a default api service
func NewNotificationApiService(c *config.Config) server.NotificationApiServicer {
	return &NotificationApiService{c}
}

// SendSlackNotification - Send a Slack message
func (s *NotificationApiService) SendSlackNotification(ctx context.Context, sendSlackMessageRequest server.SendSlackMessageRequest) (server.ImplResponse, error) {
	client := slack_lib.New(s.config.SlackAPIKey)
	timestamp, err := slack.SendMessage(sendSlackMessageRequest.To, sendSlackMessageRequest.Message, sendSlackMessageRequest.ReplyToTimestamp, client)
	if err != nil {
		fmt.Printf("Error sending slack notification: %v\n", err)
		return server.Response(http.StatusInternalServerError, nil), fmt.Errorf("Unable to send slack notification: %v", err)
	}

	return server.Response(http.StatusOK, server.SendSlackMessageResponse{Timestamp: timestamp}), nil
}
