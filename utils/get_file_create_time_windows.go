package utils

import (
	"os"
	"syscall"
	"time"
)

func GetFileLastWriteTime(path string) string {
	fileInfo, _ := os.Stat(path)
	wFileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	tNanSeconds := wFileSys.LastWriteTime.Nanoseconds()
	tSec := tNanSeconds / 1e9
	ret := time.Unix(tSec, 0).Format("2006-01-02 15:04:05")
	return ret
}
