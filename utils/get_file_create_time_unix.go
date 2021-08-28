// +build linux darwin

package utils

import (
	"os"
	"syscall"
	"time"
)

func GetFileLastWriteTime(path string) string {
	fileInfo, _ := os.Stat(path)
	statT := fileInfo.Sys().(*syscall.Stat_t)
	tCreate := statT.Mtimespec.Sec
	ret := time.Unix(tCreate, 0).Format("2006-01-02 15:04:05")
	return ret
}