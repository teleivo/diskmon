// Package usage provides a disk usage check and notification mechanism.
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

// Report is the result of a disk usage check.
type Report struct {
	Limit  uint64  // limit that the disk usage reached or exceeded
	Limits []Stats // all disk usages greater than or equal to a set limit
	Errors []error // encountered while gathering disk usages
}

// TODO embed fstat.FilesystemStat? I do not necessarily want to expose its
// HasReachedLimit method though
// Stats are disk usage statistics of a disk.
type Stats struct {
	Path  string // path on a disk that hit the set limit
	Free  uint64 // number of free bytes available
	Used  uint64 // number of used bytes
	Total uint64 // total number of bytes
}

// Notifier notifies interested parties of a usage report.
type Notifier interface {
	Notify(Report) error
}

type writeNotifier struct {
	io.Writer
}

func (wn writeNotifier) Notify(r Report) error {
	for _, l := range r.Limits {
		fmt.Fprintf(wn, "Used/Total %s/%s %q - reached limit of %d%%\n", humanize.Bytes(l.Used), humanize.Bytes(l.Total), l.Path, r.Limit)
	}
	for _, e := range r.Errors {
		fmt.Fprintf(wn, "%s\n", e.Error())
	}

	return nil
}

// WriteNotifier is a line-based usage Notifier. Every usage report stat and
// error will be printed on a dedicated line.
func WriteNotifier(w io.Writer) Notifier {
	return writeNotifier{w}
}

// Check reports disk's that reached or exceeded given limit.
func Check(basedir string, limit uint64, logger *log.Logger, nt Notifier) error {
	logger.Print("Checking disk usage")

	r, err := checkDiskUsage(basedir, limit)
	if err != nil {
		return err
	}

	if len(r.Limits) == 0 && len(r.Errors) == 0 {
		logger.Printf("Disks are below limit %d%%", limit)
		return nil
	}

	return nt.Notify(r)
}

func checkDiskUsage(basedir string, limit uint64) (Report, error) {
	r := Report{Limit: limit}
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

		if fstat.HasReachedLimit(limit) {
			r.Limits = append(r.Limits, Stats{
				Path:  file.Name(),
				Free:  fstat.Available(),
				Used:  fstat.Used(),
				Total: fstat.Total(),
			})
		}
	}
	return r, nil
}
