package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/teleivo/diskmon/fstat"
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
	err = checkUsage(*basedir, *limit, logger, out)
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
			checkUsage(*basedir, *limit, logger, out)
		}
	}
}

type notification struct {
	Limits []string
	Errors []error
}

func checkUsage(basedir string, limit uint64, logger *log.Logger, out io.Writer) error {
	logger.Print("Checking disk usage")

	n, err := checkDiskUsage(basedir, limit)
	if err != nil {
		return err
	}

	for _, l := range n.Limits {
		out.Write([]byte(l))
		out.Write([]byte("\n"))
	}
	for _, e := range n.Errors {
		out.Write([]byte(e.Error()))
		out.Write([]byte("\n"))
	}

	return nil
}

func checkDiskUsage(basedir string, limit uint64) (notification, error) {
	n := notification{}
	files, err := ioutil.ReadDir(basedir)
	if err != nil {
		return n, fmt.Errorf("error reading basedir: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		fstat, err := fstat.GetFilesystemStat(filepath.Join(basedir, file.Name()))
		if err != nil {
			n.Errors = append(n.Errors, fmt.Errorf("error getting filesystem stats from %q: %w", file.Name(), err))
			continue
		}

		if fstat.IsExceedingLimit(limit) {
			n.Limits = append(n.Limits, fmt.Sprintf("Free/Total %s/%s %q - reached limit of %d%%", humanize.Bytes(fstat.Free()), humanize.Bytes(fstat.Total()), file.Name(), limit))
		}
	}
	return n, nil
}
