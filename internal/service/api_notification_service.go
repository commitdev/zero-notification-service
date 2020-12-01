package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/commitdev/zero-notification-service/internal/server"
)

// NotificationApiService is a service that implents the logic for the NotificationApiServicer
// This service should implement the business logic for every endpoint for the NotificationApi API.
// Include any external packages or services that will be required by this service.
type NotificationApiService struct {
}

// NewNotificationApiService creates a default api service
func NewNotificationApiService() server.NotificationApiServicer {
	return &NotificationApiService{}
}

// NotificationSubscribe - Subscribe to notifications
func (s *NotificationApiService) NotificationSubscribe(ctx context.Context, subscribeRequest server.SubscribeRequest) (server.ImplResponse, error) {
	// TODO - update NotificationSubscribe with the required logic for this service method.
	// Add api_notification_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, string{}) or use other options such as http.Ok ...
	//return Response(200, string{}), nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return server.Response(http.StatusNotImplemented, nil), errors.New("NotificationSubscribe method not implemented")
}

// NotificationUnsubscribe - Unsubscribe to notifications
func (s *NotificationApiService) NotificationUnsubscribe(ctx context.Context, subscribeRequest server.SubscribeRequest) (server.ImplResponse, error) {
	// TODO - update NotificationUnsubscribe with the required logic for this service method.
	// Add api_notification_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, string{}) or use other options such as http.Ok ...
	//return Response(200, string{}), nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return server.Response(http.StatusNotImplemented, nil), errors.New("NotificationUnsubscribe method not implemented")
}

// SendNotification - Send a notification
func (s *NotificationApiService) SendNotification(ctx context.Context, notificationMessage server.NotificationMessage) (server.ImplResponse, error) {
	// TODO - update SendNotification with the required logic for this service method.
	// Add api_notification_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, string{}) or use other options such as http.Ok ...
	//return Response(200, string{}), nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return server.Response(http.StatusNotImplemented, nil), errors.New("SendNotification method not implemented")
}
