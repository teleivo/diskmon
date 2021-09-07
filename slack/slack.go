package slack

import (
	"log"
	"strings"

	"github.com/slack-go/slack"
	"github.com/teleivo/diskmon/usage"
)

type Notifier struct {
	channel string
	client  *slack.Client
	logger  *log.Logger
}

func New(token, channel string, logger *log.Logger) *Notifier {
	return &Notifier{
		channel: channel,
		client:  slack.New(token),
		logger:  logger,
	}
}

func (n *Notifier) Notify(r usage.Report) error {
	var sb strings.Builder
	l := strings.Join(r.Limits, "\n")
	if l != "" {
		sb.WriteString(l)
	}
	for _, e := range r.Errors {
		sb.WriteString(e.Error())
	}
	_, _, err := n.client.PostMessage(
		n.channel, slack.MsgOptionText(sb.String(), false),
	)
	log.Printf("Posted slack message on channel")
	return err
}
