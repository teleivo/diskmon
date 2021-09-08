package fstat

import (
	"testing"

	"golang.org/x/sys/unix"
)

func TestFilesystemStat(t *testing.T) {
	t.Run("FreeTotalUsed", func(t *testing.T) {
		fs := FilesystemStat(unix.Statfs_t{
			Ffree:  100,
			Bavail: 10,
			Blocks: 400,
			Bsize:  1024,
		})

		if want := uint64(100 * 1024); fs.Free() != want {
			t.Errorf("FilesystemStat.Free() got %d, want %d", fs.Free(), want)
		}
		if want := uint64(400 * 1024); fs.Total() != want {
			t.Errorf("FilesystemStat.Total() got %d, want %d", fs.Free(), want)
		}
		if want := uint64((400 - 100) * 1024); fs.Used() != want {
			t.Errorf("FilesystemStat.Used() got %d, want %d", fs.Free(), want)
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
					Ffree:  100,
					Bavail: 10,
					Blocks: 400,
					Bsize:  1024,
				}),
				limit: 76,
				want:  false,
			},
			"AtLimit": {
				stat: FilesystemStat(unix.Statfs_t{
					Ffree:  100,
					Bavail: 10,
					Blocks: 400,
					Bsize:  1024,
				}),
				limit: 75,
				want:  true,
			},
			"AboveLimit": {
				stat: FilesystemStat(unix.Statfs_t{
					Ffree:  100,
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
