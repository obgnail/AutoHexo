package strategy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/obgnail/AutoHexo/utils"
	"github.com/obgnail/MarkdownResouceCollecter/handler"
)

const useComments = false
const MarkdownFileHeaderTemplate = "---\ntitle: %s\ndate: %s\ntags: %s\ncomments: %t\n---\n"

// InsertHexoHeaderStrategy 插入hexo需要的header
type InsertHexoHeaderStrategy struct {
	BlogMarkdownRootDir string
}

func (s *InsertHexoHeaderStrategy) BeforeRewrite(h *handler.BaseHandler) error { return nil }

func (s *InsertHexoHeaderStrategy) AfterRewrite(h *handler.BaseHandler) error {
	return s.insertHexoHeader(h)
}

func (s *InsertHexoHeaderStrategy) insertHexoHeader(h *handler.BaseHandler) error {
	for _, f := range h.Files {
		fi, err := os.Open(f.NewPath)
		if err != nil {
			return err
		}
		title, lastWriteTime := s.getFileInfo(f.OriginPath)
		tag, err := s.getFileTag(f.NewPath)
		if err != nil {
			return err
		}
		header := s.buildHeaderContent(title, lastWriteTime, tag)
		if err := s.insertHeadIntoFile(fi, []byte(header)); err != nil {
			return err
		}
		fi.Close()
	}
	return nil
}

func (s *InsertHexoHeaderStrategy) buildHeaderContent(title, date, tag string) string {
	return fmt.Sprintf(MarkdownFileHeaderTemplate, title, date, tag, useComments)
}

func (s *InsertHexoHeaderStrategy) insertHeadIntoFile(f *os.File, header []byte) error {
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	buffer.Write(header)
	buffer.Write(content)
	if err := ioutil.WriteFile(f.Name(), buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func (s *InsertHexoHeaderStrategy) getFileInfo(filePath string) (title, LastWriteTime string) {
	fileFullPath, _, fileNameWithoutSuffix := utils.GetFilePath(filePath)

	// 除去文件名的`书籍__`前缀
	prefixAndFileName := strings.Split(fileNameWithoutSuffix, "__")
	if len(prefixAndFileName) > 1 {
		title = prefixAndFileName[1]
	} else {
		title = fileNameWithoutSuffix
	}
	LastWriteTime = utils.GetFileLastWriteTime(fileFullPath)
	return
}

// fileFullPath:        /Users/heyingliang/myTemp/blog/source/_posts/Binlog/MySQL Binlog 介绍.md
// BlogMarkdownRootDir: /Users/heyingliang/myTemp/blog/source/_posts
// want: Binlog
func (s *InsertHexoHeaderStrategy) getFileTag(fileFullPath string) (tag string, err error) {
	path, err := filepath.Rel(s.BlogMarkdownRootDir, fileFullPath)
	if err != nil {
		return
	}
	ss := strings.Split(path, string(os.PathSeparator))
	if len(ss) >= 2 {
		tag = ss[0]
	}
	return
}
