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

	msg := n.Message(r, host)

	// TODO should we use a context? What is the default behavior, does posting
	// the message ever time out?
	_, _, err = n.client.PostMessage(n.channel, msg)
	log.Printf("Posted slack message on channel")

	return err
}

func (n *Notifier) Message(r usage.Report, host string) slack.MsgOption {
	// TODO clean up code
	header := "Disk usage report"
	if host != "" {
		header = header + fmt.Sprintf(" for host %q", host)
	}
	headerText := slack.NewTextBlockObject("plain_text", header, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	limitImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/notificationsWarningIcon.png", "notifications warning icon")
	limitHeader := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Following disks have reached the usage limit of %d%%*", r.Limits[0].Limit), false, false)

	var sb strings.Builder
	for _, l := range r.Limits {
		fmt.Fprintf(&sb, "• %q - %s/%s (free/total)\n", l.Path, humanize.Bytes(l.Free), humanize.Bytes(l.Total))
	}
	// TODO split errors into own section?
	for _, e := range r.Errors {
		sb.WriteString(e.Error())
	}
	rest := slack.NewTextBlockObject("mrkdwn", sb.String(), false, false)
	restSection := slack.NewSectionBlock(rest, nil, nil)

	limitSection := slack.NewContextBlock(
		"",
		[]slack.MixedElement{limitImage, limitHeader}...,
	)
	divSection := slack.NewDividerBlock()

	return slack.MsgOptionBlocks(headerSection, divSection, limitSection, restSection)
}
