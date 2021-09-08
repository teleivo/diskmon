package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/teleivo/diskmon/slack"
	"github.com/teleivo/diskmon/usage"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprint(os.Stderr, err, "\n")
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	basedir := flags.String("basedir", "", "statfs syscall information will be gathered for each directory (depth 1) in this base directory")
	limit := flags.Uint64("limit", 80, "notification will be sent if a disk's usage is greater than or equal to given limit in percentage")
	interval := flags.Uint("interval", 60, "interval in minutes at which the disk usage will be checked")
	slackToken := flags.String("slackToken", "", "Slack Bot User OAuth Token used to post notifications to Slack")
	slackChannel := flags.String("slackChannel", "", "Slack channel ID where notifications are posted to")
	err := flags.Parse(args[1:])
	if err != nil {
		return err
	}
	if *basedir == "" {
		return errors.New("basedir must be provided")
	}
	if *slackToken != "" && *slackChannel == "" || *slackChannel != "" && *slackToken == "" {
		return errors.New("both slackChannel and slackApiToken must be either provided or not")
	}

	logger := log.New(out, args[0]+" ", log.LUTC)
	var notifier usage.Notifier
	if *slackToken != "" && *slackChannel != "" {
		notifier = slack.New(*slackToken, *slackChannel, logger)
	} else {
		notifier = usage.WriteNotifier(out)
	}

	t := time.NewTicker(time.Minute * time.Duration(*interval))
	defer t.Stop()

	// check usage once right after starting up
	err = usage.Check(*basedir, *limit, logger, notifier)
	if err != nil {
		return err
	}
	for {
		select {
		case <-t.C:
			err := usage.Check(*basedir, *limit, logger, notifier)
			// TODO what if there is an error here on reading the basedir? At
			// least log it for now
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}
