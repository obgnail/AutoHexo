package utils

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func GetFilePath(f *os.File) (fileFullPath, fileNameWithSuffix, fileNameWithoutSuffix string) {
	fileFullPath = f.Name()
	fileNameWithSuffix = path.Base(fileFullPath)
	fileSuffix := path.Ext(fileNameWithSuffix)
	fileNameWithoutSuffix = strings.TrimSuffix(fileNameWithSuffix, fileSuffix)
	return
}

func GetFileInfo(f *os.File) (title, date, tags string, err error) {
	fmt.Println(f)
	fileFullPath, _, fileNameWithoutSuffix := GetFilePath(f)

	// 除去文件名的`书籍__`前缀
	fileName := strings.Split(fileNameWithoutSuffix, "__")
	if len(fileName) >= 0 {
		title = fileName[1]
	}

	d := GetFileCreateTime(fileFullPath)
	date = fmt.Sprintf("%d", d)

	learningDir := path.Join(markdownDirPath,"Learning")
}

//
//func InsertContentIntoFileHead(filePath string, content []byte) error {
//	originContent, _ := ReadAll(filePath)
//	var buffer bytes.Buffer
//	buffer.Write(content)
//	buffer.Write(originContent)
//	if err := ioutil.WriteFile(filePath, buffer.Bytes(), 0644); err != nil {
//		return err
//	}
//	return nil
//}
//

func GetFileCreateTime(path string) int64 {
	osType := runtime.GOOS
	fileInfo, _ := os.Stat(path)
	if osType == "windows" {
		wFileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
		tNanSeconds := wFileSys.CreationTime.Nanoseconds()
		tSec := tNanSeconds / 1e9
		return tSec
	}
	//if osType == "linux" {
	//	stat_t := fileInfo.Sys().(*syscall.Stat_t)
	//	tCreate := int64(stat_t.Ctim.Sec)
	//	return tCreate
	//}
	return time.Now().Unix()
}
