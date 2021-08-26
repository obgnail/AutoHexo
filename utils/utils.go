package utils

import (
	"os"
	"path"
	"strings"
)

func GetFilePath(f *os.File) (fileFullPath, fileNameWithSuffix, fileNameWithoutSuffix string) {
	fileFullPath = f.Name()
	fileNameWithSuffix = path.Base(fileFullPath)
	fileSuffix := path.Ext(fileNameWithSuffix)
	fileNameWithoutSuffix = strings.TrimSuffix(fileNameWithSuffix, fileSuffix)
	return
}
