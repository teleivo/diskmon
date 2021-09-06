package fstat

import (
	"golang.org/x/sys/unix"
)

type FilesystemStat unix.Statfs_t

func (fs FilesystemStat) Used() uint64 {
	return (fs.Blocks - fs.Bavail) * uint64(fs.Bsize)
}

func (fs FilesystemStat) Free() uint64 {
	return fs.Bavail * uint64(fs.Bsize)
}

func (fs FilesystemStat) Total() uint64 {
	return fs.Blocks * uint64(fs.Bsize)
}

func (fs FilesystemStat) IsExceedingLimit(limit uint64) bool {
	return uint64((float64(fs.Used())/float64(fs.Total()))*100) > limit
}

func GetFilesystemStat(path string) (FilesystemStat, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return FilesystemStat{}, err
	}

	return FilesystemStat(stat), nil
}
