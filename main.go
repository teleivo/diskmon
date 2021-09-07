package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

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
	basedir := flags.String("basedir", "", "statfs syscall information will be printed for each directory (depth 1) in this base directory")
	limit := flags.Uint64("limit", 80, "percentage of disk usage at which notification should be sent")
	interval := flags.Uint("interval", 60, "interval in minutes at which the disk usage will be checked")
	err := flags.Parse(args[1:])
	if err != nil {
		return err
	}
	if *basedir == "" {
		return errors.New("basedir must be provided")
	}

	t := time.NewTicker(time.Minute * time.Duration(*interval))
	defer t.Stop()

	logger := log.New(out, args[0]+" ", log.LUTC)

	// check usage once right after starting up
	err = usage.Check(*basedir, *limit, logger, out)
	if err != nil {
		return err
	}
	for {
		select {
		case <-t.C:
			// TODO what if there is an error here on reading the basedir?
			// I think its a good idea to call ReadDir every time we check
			// usage since one could add a volume to a droplet after the
			// diskmon has started
			usage.Check(*basedir, *limit, logger, out)
		}
	}
}
