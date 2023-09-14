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
	cl.Called(email)
	return nil, nil
}

func TestSendMail(t *testing.T) {
	toList := createRandomRecipients(1, 1)
	cc := make([]server.EmailRecipient, 0)
	bcc := make([]server.EmailRecipient, 0)
	from := server.EmailSender{Name: "Test User", Address: "address@domain.com"}
	headers := map[string]string{
		"X-Test-Header": "Test Header Value",
	}
	message := server.MailMessage{Subject: "Subject", Body: "Body"}
	client := FakeClient{}

	headersMatcher := mock.MatchedBy(func(m *sendgridMail.SGMailV3) bool {
		return m.Headers["X-Test-Header"] == headers["X-Test-Header"]
	})
	client.On("Send", headersMatcher).Return(nil, nil)
	mail.SendIndividualMail(toList, from, cc, bcc, headers, message, &client, []string{}, map[string]string{})
	client.AssertNumberOfCalls(t, "Send", 1)
}

func TestSendBulkMail(t *testing.T) {
	toList := createRandomRecipients(2, 5)
	cc := make([]server.EmailRecipient, 0)
	bcc := make([]server.EmailRecipient, 0)
	from := server.EmailSender{Name: "Test User", Address: "address@domain.com"}
	headers := make(map[string]string)
	message := server.MailMessage{Subject: "Subject", Body: "Body"}
	client := FakeClient{}

	client.On("Send", mock.Anything).Return(nil, nil)

	responseChannel := make(chan mail.BulkSendAttempt)
	mail.SendBulkMail(toList, from, cc, bcc, headers, message, &client, responseChannel)

	// Range over the channel until empty
	returnedCount := 0
	for range responseChannel {
		returnedCount++
	}

	assert.Equal(t, len(toList), returnedCount, "Response count should match requests sent")

	// Check that the send function was called
	client.AssertNumberOfCalls(t, "Send", len(toList))
}

func TestRemoveInvalidRecipients(t *testing.T) {
	toList := createRandomRecipients(2, 5)

	originalSize := len(toList)

	toList[0].Address = "address@commit.dev"

	alteredList := mail.RemoveInvalidRecipients(toList, []string{"commit.dev", "domain.com"})
	assert.Equal(t, len(alteredList), originalSize, "All addresses should remain in the list")

	alteredList = mail.RemoveInvalidRecipients(toList, []string{"commit.dev"})
	assert.Equal(t, len(alteredList), 1, "1 address should remain in the list")
}

// createRandomRecipients creates a random list of recipients
func createRandomRecipients(min int, randCount int) []server.EmailRecipient {
	var toList []server.EmailRecipient
	// Create a random number of mails
	rand.Seed(time.Now().UnixNano())
	sendCount := rand.Intn(randCount) + min
	for i := 0; i < sendCount; i++ {
		toList = append(toList, server.EmailRecipient{
			Name:    fmt.Sprintf("Test Recipient %d", i),
			Address: fmt.Sprintf("address%d@domain.com", i),
		})
	}
	return toList
}
