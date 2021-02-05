package slack

import (
	"github.com/commitdev/zero-notification-service/internal/server"
	slack_lib "github.com/slack-go/slack"
)

type Client interface {
	PostMessage(channelID string, options ...slack_lib.MsgOption) (string, string, error)
}

// SendIndividualMail sends an email message
func SendMessage(to server.SlackRecipient, message server.SlackMessage, replyToTimestamp string, client Client) (string, error) {
	args := []slack_lib.MsgOption{
		slack_lib.MsgOptionText(message.Body, false),
		slack_lib.MsgOptionAsUser(true),
	}
	if replyToTimestamp != "" {
		args = append(args, slack_lib.MsgOptionTS(replyToTimestamp))
	}

	_, timestamp, err := client.PostMessage(to.ConversationId, args...)
	return timestamp, err
}
