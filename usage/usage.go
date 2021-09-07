package usage

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/teleivo/diskmon/fstat"
)

type Notification struct {
	Limits []string
	Errors []error
}

func Check(basedir string, limit uint64, logger *log.Logger, out io.Writer) error {
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

func checkDiskUsage(basedir string, limit uint64) (Notification, error) {
	n := Notification{}
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
