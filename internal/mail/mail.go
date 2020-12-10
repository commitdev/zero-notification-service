package mail

import (
	"sync"

	"github.com/commitdev/zero-notification-service/internal/server"
	"github.com/sendgrid/rest"
	sendgridMail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type BulkSendAttempt struct {
	EmailAddress string
	Response     *rest.Response
	Error        error
}

type Client interface {
	Send(email *sendgridMail.SGMailV3) (*rest.Response, error)
}

// SendBulkMail sends a batch of email messages to all the specified recipients
// All the calls to send mail happen in parallel, with their responses returned on the provided channel
func SendBulkMail(toList []server.Recipient, from server.Sender, message server.MailMessage, client Client, responseChannel chan BulkSendAttempt) {
	wg := sync.WaitGroup{}
	wg.Add(len(toList))

	// Create goroutines for each send
	for _, to := range toList {
		go func(to server.Recipient) {
			response, err := SendIndividualMail(to, from, message, client)
			responseChannel <- BulkSendAttempt{to.Address, response, err}
			wg.Done()
		}(to)
	}
	// Wait on the all responses to close the channel to signal that the operation is complete
	go func() {
		wg.Wait()
		close(responseChannel)
	}()
}

// SendIndividualMail sends an email message
func SendIndividualMail(to server.Recipient, from server.Sender, message server.MailMessage, client Client) (*rest.Response, error) {
	sendFrom := sendgridMail.NewEmail(from.Name, from.Address)
	sendTo := sendgridMail.NewEmail(to.Name, to.Address)
	sendMessage := sendgridMail.NewSingleEmail(sendFrom, message.Subject, sendTo, message.Body, message.RichBody)
	sendMessage.SetTemplateID(message.TemplateId)
	sendMessage.SetSendAt(int(message.ScheduleSendAtTimestamp))
	return client.Send(sendMessage)
}
