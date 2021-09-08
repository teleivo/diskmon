package slack

import (
	"fmt"
	"log"
	"strings"

	"github.com/dustin/go-humanize"
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
	for _, l := range r.Limits {
		fmt.Fprintf(&sb, "Free/Total %s/%s %q - reached limit of %d%%\n", humanize.Bytes(l.Free), humanize.Bytes(l.Total), l.Path, l.Limit)
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
