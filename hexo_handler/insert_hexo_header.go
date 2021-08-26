package hexo_handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/obgnail/AutoHexo/utils"
	"github.com/obgnail/MarkdownResouceCollecter/handler"
)

const MarkdownFileHeaderTemplate = "---\ntitle: %s\ndate: %s\ntags: %s\n---\n"

// InsertHexoHeaderStrategy 插入hexo需要的header
type InsertHexoHeaderStrategy struct {
	blogMarkdownRootDir string
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
		title, createDate := s.getFileInfo(fi)
		tag, err := s.getFileTag(f.NewPath)
		if err != nil {
			return err
		}
		header := s.buildHeaderContent(title, createDate, tag)
		if err := s.insertHeadIntoFile(fi, []byte(header)); err != nil {
			return err
		}
		fi.Close()
	}
	return nil
}

func (s *InsertHexoHeaderStrategy) buildHeaderContent(title, date, tag string) string {
	return fmt.Sprintf(MarkdownFileHeaderTemplate, title, date, tag)
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

func (s *InsertHexoHeaderStrategy) getFileInfo(f *os.File) (title, createDate string) {
	fileFullPath, _, fileNameWithoutSuffix := utils.GetFilePath(f)

	// 除去文件名的`书籍__`前缀
	prefixAndFileName := strings.Split(fileNameWithoutSuffix, "__")
	if len(prefixAndFileName) > 1 {
		title = prefixAndFileName[1]
	} else {
		title = fileNameWithoutSuffix
	}
	createDate = utils.GetFileCreateTime(fileFullPath)
	return
}

// fileFullPath:        /Users/heyingliang/myTemp/blog/source/_posts/Binlog/MySQL Binlog 介绍.md
// blogMarkdownRootDir: /Users/heyingliang/myTemp/blog/source/_posts
// want: Binlog
func (s *InsertHexoHeaderStrategy) getFileTag(fileFullPath string) (tag string, err error) {
	regexpString := filepath.Join(s.blogMarkdownRootDir, `(.+?)`) + string(os.PathSeparator)
	re, err := regexp.Compile(regexpString)
	if err != nil {
		return
	}
	matches := re.FindAllStringSubmatch(fileFullPath, -1)
	if matches == nil {
		return "", fmt.Errorf("[Error] get file tag fail : %s", fileFullPath)
	}
	tag = matches[0][1]
	return
}
