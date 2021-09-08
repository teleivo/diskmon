package fstat

import (
	"golang.org/x/sys/unix"
)

type FilesystemStat unix.Statfs_t

// Used returns the number of used bytes.
func (fs FilesystemStat) Used() uint64 {
	return (fs.Blocks - fs.Ffree) * uint64(fs.Bsize)
}

// Free returns the number of free bytes available. Note that a non-privileged
// user might have less bytes available (see Statfs_t.Bavail).
func (fs FilesystemStat) Free() uint64 {
	return fs.Ffree * uint64(fs.Bsize)
}

// Total returns the total number of bytes.
func (fs FilesystemStat) Total() uint64 {
	return fs.Blocks * uint64(fs.Bsize)
}

// HasReachedLimit returns true if the used disk space is greater than or equal
// to the given limit in percent.
func (fs FilesystemStat) HasReachedLimit(limit uint64) bool {
	used := uint64((float64(fs.Used()) / float64(fs.Total())) * 100)
	return used >= limit
}

func GetFilesystemStat(path string) (FilesystemStat, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return FilesystemStat{}, err
	}

	return FilesystemStat(stat), nil
}
