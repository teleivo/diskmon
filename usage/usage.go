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

type Report struct {
	Limits []string
	Errors []error
}

type Notifier interface {
	Notify(Report) error
}

type writeNotifier struct {
	io.Writer
}

func (wn writeNotifier) Notify(r Report) error {
	for _, l := range r.Limits {
		wn.Write([]byte(l))
		wn.Write([]byte("\n"))
	}
	for _, e := range r.Errors {
		wn.Write([]byte(e.Error()))
		wn.Write([]byte("\n"))
	}

	return nil
}

func WriteNotifier(w io.Writer) Notifier {
	return writeNotifier{w}
}

func Check(basedir string, limit uint64, logger *log.Logger, nt Notifier) error {
	logger.Print("Checking disk usage")

	r, err := checkDiskUsage(basedir, limit)
	if err != nil {
		return err
	}

	if len(r.Limits) == 0 && len(r.Errors) == 0 {
		logger.Print("No limits exceeded and no errors found")
		return nil
	}

	return nt.Notify(r)
}

func checkDiskUsage(basedir string, limit uint64) (Report, error) {
	r := Report{}
	files, err := ioutil.ReadDir(basedir)
	if err != nil {
		return r, fmt.Errorf("error reading basedir: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		fstat, err := fstat.GetFilesystemStat(filepath.Join(basedir, file.Name()))
		if err != nil {
			r.Errors = append(r.Errors, fmt.Errorf("error getting filesystem stats from %q: %w", file.Name(), err))
			continue
		}

		if fstat.IsExceedingLimit(limit) {
			r.Limits = append(r.Limits, fmt.Sprintf("Free/Total %s/%s %q - reached limit of %d%%", humanize.Bytes(fstat.Free()), humanize.Bytes(fstat.Total()), file.Name(), limit))
		}
	}
	return r, nil
}
