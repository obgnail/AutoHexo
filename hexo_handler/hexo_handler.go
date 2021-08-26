package hexo_handler

import (
	handlerConfig "github.com/obgnail/MarkdownResouceCollecter/config"
	"github.com/obgnail/MarkdownResouceCollecter/handler"
	"log"
	"strings"
	"sync"
)

// HexoCreateHandlers 用于OnCreate事件
type HexoCreateHandlers struct {
	markdownRootDir        string
	blogResourceRootDir    string
	blogMarkdownRootDir    string
	markdownFileSuffix     string
	localPictureUseAbsPath bool

	handlers []*handler.BaseHandler
}

func New(
	markdownRootDir, blogResourceRootDir, blogMarkdownRootDir string,
	markdownFileSuffix string, localPictureUseAbsPath bool, filesPath []string,
) *HexoCreateHandlers {
	hh := &HexoCreateHandlers{
		markdownRootDir:        markdownRootDir,
		blogResourceRootDir:    blogResourceRootDir,
		blogMarkdownRootDir:    blogMarkdownRootDir,
		markdownFileSuffix:     markdownFileSuffix,
		localPictureUseAbsPath: localPictureUseAbsPath,
	}

	for _, fp := range filesPath {
		blogMarkdownPath := hh.GetBlogMarkdownPath(fp)
		h := NewHandler(fp, hh.blogMarkdownRootDir, hh.blogResourceRootDir, blogMarkdownPath, hh.markdownFileSuffix, hh.localPictureUseAbsPath)
		hh.AppendHandler(h)
	}
	return hh
}

func NewHandler(
	originMarkdownRootPath, blogMarkdownRootDirPath, blogResourceRootDirPath, blogMarkdownFilePath string,
	markdownFileSuffix string, localPictureUseAbsPath bool,
) *handler.BaseHandler {
	cfg := handlerConfig.InitConfig(
		originMarkdownRootPath, blogResourceRootDirPath, blogMarkdownFilePath,
		nil, markdownFileSuffix, localPictureUseAbsPath,
	)
	h := handler.New(cfg)
	h.AppendStrategy(&handler.CollectNetWorkPictureStrategy{})
	h.AppendStrategy(&handler.CollectLocalPictureStrategy{})
	h.AppendStrategy(&InsertHexoHeaderStrategy{blogMarkdownRootDirPath})
	return h
}

func (hh *HexoCreateHandlers) AppendHandler(h *handler.BaseHandler) {
	hh.handlers = append(hh.handlers, h)
}

// input:  /Users/heyingliang/Dropbox/root/md/Learning/Canal/网站__Cannal入门
// output: /Users/heyingliang/myTemp/root2/md4/Learning/Canal/网站__Cannal入门
func (hh *HexoCreateHandlers) GetBlogMarkdownPath(originMarkdownDir string) string {
	return strings.Replace(originMarkdownDir, hh.markdownRootDir, hh.blogMarkdownRootDir, 1)
}

func (hh *HexoCreateHandlers) Run() {
	var wg sync.WaitGroup
	for _, h := range hh.handlers {
		wg.Add(1)
		go func(h *handler.BaseHandler) {
			defer wg.Done()
			if err := h.Run(); err != nil {
				log.Println("[Error] handler run err:", err)
			}
		}(h)
	}
	wg.Wait()
}
