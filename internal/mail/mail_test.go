package mail_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/commitdev/zero-notification-service/internal/mail"
	"github.com/commitdev/zero-notification-service/internal/server"
	"github.com/sendgrid/rest"
	sendgridMail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FakeClient struct {
	mock.Mock
}

// Mock the Send function
func (cl *FakeClient) Send(email *sendgridMail.SGMailV3) (*rest.Response, error) {
	cl.Called()
	return nil, nil
}

func TestSendBulkMail(t *testing.T) {
	var toList []server.Recipient
	// Create a random number of mails
	rand.Seed(time.Now().UnixNano())
	sendCount := rand.Intn(5) + 2
	for i := 0; i < sendCount; i++ {
		toList = append(toList, server.Recipient{fmt.Sprintf("Test Recipient %d", i), fmt.Sprintf("address%d@domain.com", i)})
	}
	from := server.Sender{"Test User", "address@domain.com"}
	message := server.MailMessage{Subject: "Subject", Body: "Body"}
	client := FakeClient{}

	client.On("Send").Return(nil, nil)

	responseChannel := make(chan mail.BulkSendAttempt)
	mail.SendBulkMail(toList, from, message, &client, responseChannel)

	// Range over the channel until empty
	returnedCount := 0
	for range responseChannel {
		returnedCount++
	}

	assert.Equal(t, sendCount, returnedCount, "Response count should match requests sent")

	// Check that the send function was called
	client.AssertNumberOfCalls(t, "Send", sendCount)
}
