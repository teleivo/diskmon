package fstat

import (
	"testing"

	"golang.org/x/sys/unix"
)

func TestFilesystemStat(t *testing.T) {
	t.Run("IsExceedingLimit", func(t *testing.T) {
		// TODO use Bavail or Bfree?
		fs := FilesystemStat(unix.Statfs_t{
			Bavail: 100,
			Blocks: 400,
			Bsize:  1024,
		})

		if want := uint64(100 * 1024); fs.Free() != want {
			t.Errorf("FilesystemStat.Free got %d, want %d", fs.Free(), want)
		}
		if want := uint64(400 * 1024); fs.Total() != want {
			t.Errorf("FilesystemStat.Total got %d, want %d", fs.Free(), want)
		}
		if want := uint64((400 - 100) * 1024); fs.Used() != want {
			t.Errorf("FilesystemStat.Used got %d, want %d", fs.Free(), want)
		}

		if fs.IsExceedingLimit(74) != true {
			t.Errorf("IsExceedingLimit(74) should be true")
		}

		if fs.IsExceedingLimit(75) != false {
			t.Errorf("IsExceedingLimit(75) should be false")
		}

		if fs.IsExceedingLimit(76) != false {
			t.Errorf("IsExceedingLimit(76) should be false")
		}
	})
}
