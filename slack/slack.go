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

type notifier struct {
	channel string
	client  *slack.Client
	logger  *log.Logger
}

func New(token, channel string, logger *log.Logger) notifier {
	return notifier{
		channel: channel,
		client:  slack.New(token),
		logger:  logger,
	}
}

func (n notifier) Notify(r usage.Report) error {
	host, err := os.Hostname()
	if err != nil {
		n.logger.Printf("Failed to get hostname %v", err)
	}

	// TODO should we use a context? What is the default behavior, does posting
	// the message ever time out?
	_, _, err = n.client.PostMessage(
		n.channel,
		slack.MsgOptionBlocks(formatMessage(r, host)...),
	)
	// TODO I might not have posted a message on slack. Think about how to
	// handle the error
	n.logger.Printf("Posted slack message on channel")

	return err
}

func formatMessage(r usage.Report, host string) []slack.Block {
	header := "Disk usage report"
	if host != "" {
		header = header + fmt.Sprintf(" for host %q", host)
	}
	headerText := slack.NewTextBlockObject("plain_text", header, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	limitHeader := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("⚠️ *Following disks have reached the usage limit of %d%%*", r.Limit), false, false)
	var sb strings.Builder
	for _, l := range r.Limits {
		perc := uint64((float64(l.Used) / float64(l.Total)) * 100)
		fmt.Fprintf(&sb, "• %q - *%d%% full* - %s/%s (used/total)\n", l.Path, perc, humanize.Bytes(l.Used), humanize.Bytes(l.Total))
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
