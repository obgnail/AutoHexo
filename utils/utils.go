package utils

import (
	"path/filepath"
	"strings"
)

func GetFilePath(filePath string) (fileFullPath, fileNameWithSuffix, fileNameWithoutSuffix string) {
	fileFullPath = filePath
	fileNameWithSuffix = filepath.Base(fileFullPath)
	fileSuffix := filepath.Ext(fileNameWithSuffix)
	fileNameWithoutSuffix = strings.TrimSuffix(fileNameWithSuffix, fileSuffix)
	return
}
