package fstat

import (
	"testing"

	"golang.org/x/sys/unix"
)

func TestFilesystemStat(t *testing.T) {
	t.Run("AvailableTotalUsed", func(t *testing.T) {
		fs := FilesystemStat(unix.Statfs_t{
			Ffree:  200, // total free disk space
			Bavail: 100, // free disk space accessible by a non-privileged user
			Blocks: 400,
			Bsize:  1024,
		})

		if want := uint64(100 * 1024); fs.Available() != want {
			t.Errorf("FilesystemStat.Available() got %d, want %d", fs.Available(), want)
		}
		if want := uint64(400 * 1024); fs.Total() != want {
			t.Errorf("FilesystemStat.Total() got %d, want %d", fs.Total(), want)
		}
		if want := uint64((400 - 100) * 1024); fs.Used() != want {
			t.Errorf("FilesystemStat.Used() got %d, want %d", fs.Used(), want)
		}
	})
	t.Run("HasReachedLimit", func(t *testing.T) {
		tt := map[string]struct {
			stat  FilesystemStat
			limit uint64
			want  bool
		}{
			"BelowLimit": {
				stat: FilesystemStat(unix.Statfs_t{
					Ffree:  200,
					Bavail: 100,
					Blocks: 400,
					Bsize:  1024,
				}),
				limit: 76,
				want:  false,
			},
			"AtLimit": {
				stat: FilesystemStat(unix.Statfs_t{
					Ffree:  200,
					Bavail: 100,
					Blocks: 400,
					Bsize:  1024,
				}),
				limit: 75,
				want:  true,
			},
			"AboveLimit": {
				stat: FilesystemStat(unix.Statfs_t{
					Ffree:  200,
					Bavail: 100,
					Blocks: 400,
					Bsize:  1024,
				}),
				limit: 74,
				want:  true,
			},
		}

		for s, tc := range tt {
			t.Run(s, func(t *testing.T) {
				if got := tc.stat.HasReachedLimit(tc.limit); got != tc.want {
					t.Errorf("got %t, want %t", got, tc.want)
				}
			})
		}
	})
}
