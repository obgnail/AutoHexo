package utils

import (
	"os"
	"syscall"
	"time"
)

func GetFileCreateTime(path string) string {
	fileInfo, _ := os.Stat(path)
	statT := fileInfo.Sys().(*syscall.Stat_t)
	tCreate := statT.Ctimespec.Sec
	ret := time.Unix(tCreate, 0).Format("2006-01-02 15:04:05")
	return ret
}