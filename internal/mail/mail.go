package mail

import (
	"fmt"
	"strings"
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
func SendBulkMail(toList []server.EmailRecipient, from server.EmailSender, cc []server.EmailRecipient, bcc []server.EmailRecipient, headers map[string]string, message server.MailMessage, client Client, responseChannel chan BulkSendAttempt) {
	wg := sync.WaitGroup{}
	wg.Add(len(toList))

	// Create goroutines for each send
	for _, to := range toList {
		go func(to server.EmailRecipient) {
			response, err := SendIndividualMail([]server.EmailRecipient{to}, from, cc, bcc, headers, message, client, []string{}, map[string]string{})
			responseChannel <- BulkSendAttempt{to.Address, response, err}
			wg.Done()
		}(to)
	}
	// Wait on all the responses to close the channel to signal that the operation is complete
	go func() {
		wg.Wait()
		close(responseChannel)
	}()
}

// SendIndividualMail sends an email message
func SendIndividualMail(to []server.EmailRecipient, from server.EmailSender, cc []server.EmailRecipient, bcc []server.EmailRecipient, headers map[string]string, message server.MailMessage, client Client, tags []string, customArgs map[string]string) (*rest.Response, error) {
	sendMessage := sendgridMail.NewV3Mail()

	sendMessage.SetFrom(sendgridMail.NewEmail(from.Name, from.Address))

	if message.Body != "" {
		sendMessage.AddContent(sendgridMail.NewContent("text/plain", message.Body))
	}
	if message.RichBody != "" {
		sendMessage.AddContent(sendgridMail.NewContent("text/html", message.RichBody))
	}

	sendMessage.SetTemplateID(message.TemplateId)
	sendMessage.SetSendAt(int(message.ScheduleSendAtTimestamp))

	personalization := sendgridMail.NewPersonalization()

	personalization.Subject = message.Subject

	personalization.AddTos(convertAddresses(to)...)

	if len(cc) > 0 {
		personalization.AddCCs(convertAddresses(cc)...)
	}

	if len(bcc) > 0 {
		personalization.AddBCCs(convertAddresses(bcc)...)
	}
	sendMessage.AddPersonalizations(personalization)

	if len(headers) > 0 {
		for key, value := range headers {
			sendMessage.SetHeader(key, value)
		}
	}

	if len(tags) > 0 {
		sendMessage.AddCategories(tags...)
	}

	if len(customArgs) > 0 {
		for key, value := range customArgs {
			sendMessage.SetCustomArg(key, value)
		}
	}

	return client.Send(sendMessage)
}

// convertAddresses converts a list of EmailRecipient structs to a list of sendgrid's email address type
func convertAddresses(addresses []server.EmailRecipient) []*sendgridMail.Email {
	returnAddresses := make([]*sendgridMail.Email, len(addresses))
	for i, address := range addresses {
		returnAddresses[i] = sendgridMail.NewEmail(address.Name, address.Address)
	}
	return returnAddresses
}

// RemoveInvalidRecipients accepts a list of recipients and removes the ones with domains not in the allowed list
func RemoveInvalidRecipients(recipients []server.EmailRecipient, allowedDomains []string) []server.EmailRecipient {
	valid := []server.EmailRecipient{}
	for _, recipient := range recipients {
		if addressInAllowedDomain(recipient.Address, allowedDomains) {
			valid = append(valid, recipient)
		}
	}
	return valid
}

// addressInAllowedDomain checks if a single email address is in a list of domains
func addressInAllowedDomain(address string, domains []string) bool {
	for _, domain := range domains {
		if strings.HasSuffix(address, fmt.Sprintf("@%s", domain)) {
			return true
		}
	}
	return false
}
