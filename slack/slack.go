package slack

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/slack-go/slack"
	"github.com/teleivo/diskmon/usage"
)

// TODO make private
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
	host, err := os.Hostname()
	if err != nil {
		log.Printf("Failed to get hostname %v", err)
	}

	msg := n.message(r, host)

	// TODO should we use a context? What is the default behavior, does posting
	// the message ever time out?
	_, _, err = n.client.PostMessage(n.channel, msg)
	log.Printf("Posted slack message on channel")

	return err
}

func (n *Notifier) message(r usage.Report, host string) slack.MsgOption {
	return slack.MsgOptionBlocks(formatMessage(r, host)...)
}

func formatMessage(r usage.Report, host string) []slack.Block {
	// TODO move limit to usage.Report.Limit instead of with the usage.Stat
	header := "Disk usage report"
	if host != "" {
		header = header + fmt.Sprintf(" for host %q", host)
	}
	headerText := slack.NewTextBlockObject("plain_text", header, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	limitHeader := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("⚠️ *Following disks have reached the usage limit of %d%%*", r.Limits[0].Limit), false, false)
	var sb strings.Builder
	for _, l := range r.Limits {
		fmt.Fprintf(&sb, "• %q - %s/%s (free/total)\n", l.Path, humanize.Bytes(l.Free), humanize.Bytes(l.Total))
	}
	limits := slack.NewTextBlockObject("mrkdwn", sb.String(), false, false)
	limitSection := slack.NewContextBlock("limits", limitHeader, limits)

	divSection := slack.NewDividerBlock()
	msg := []slack.Block{slack.Block(headerSection), divSection, limitSection}

	if len(r.Errors) > 0 {
		sb.Reset()
		errorHeader := slack.NewTextBlockObject("mrkdwn", "⛔️ *Following error(s) have been encountered while gathering disk usage stats*", false, false)
		for _, e := range r.Errors {
			fmt.Fprintf(&sb, "• %s\n", e.Error())
		}
		errors := slack.NewTextBlockObject("mrkdwn", sb.String(), false, false)
		errorsSection := slack.NewContextBlock("errors", errorHeader, errors)
		msg = append(msg, errorsSection)
	}

	return msg
}
