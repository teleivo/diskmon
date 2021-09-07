package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
	// check usage once right after starting up
	checkUsage(*basedir, *limit, out)
	for {
		select {
		case <-t.C:
			checkUsage(*basedir, *limit, out)
		}
	}
}

func checkUsage(basedir string, limit uint64, out io.Writer) error {
	fmt.Fprintf(out, "Checking disk usage\n")
	files, err := ioutil.ReadDir(basedir)
	if err != nil {
		return fmt.Errorf("error reading basedir: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		fstat, err := fstat.GetFilesystemStat(filepath.Join(basedir, file.Name()))
		if err != nil {
			// TODO what do to with such errors? also send a notification?
			err = fmt.Errorf("error getting filesystem stats from %q: %w", file.Name(), err)
			// TODO remove printing a newline once I have decided on how to handle errors
			fmt.Fprint(out, err, "\n")
			continue
		}

		if fstat.IsExceedingLimit(limit) {
			fmt.Fprintf(out, "Free/Total %s/%s %q - reached limit of %d%%\n", humanize.Bytes(fstat.Free()), humanize.Bytes(fstat.Total()), file.Name(), limit)
		}
	}
	return nil
}
