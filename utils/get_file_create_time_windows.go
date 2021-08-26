package utils

import (
	"os"
	"syscall"
)

func GetFileCreateTime(path string) int64 {
	fileInfo, _ := os.Stat(path)
	wFileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	tNanSeconds := wFileSys.CreationTime.Nanoseconds()
	tSec := tNanSeconds / 1e9
	return tSec
}
