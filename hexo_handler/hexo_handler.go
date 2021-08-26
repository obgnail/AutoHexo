package hexo_handler

import (
	"strings"

	handlerConfig "github.com/obgnail/MarkdownResouceCollecter/config"
	"github.com/obgnail/MarkdownResouceCollecter/handler"
)

// 用于OnCreate事件
func New(
	originMarkdownRootDir, blogResourceRootDir, blogMarkdownRootDir string,
	markdownFileSuffix string, localPictureUseAbsPath bool, changedFilePath string,
) *handler.BaseHandler {
	blogMarkdownPath := strings.Replace(changedFilePath, originMarkdownRootDir, blogMarkdownRootDir, 1)
	cfg := handlerConfig.InitConfig(
		changedFilePath, blogResourceRootDir, blogMarkdownPath,
		nil, markdownFileSuffix, localPictureUseAbsPath,
	)
	h := handler.New(cfg)
	h.AppendStrategy(&handler.CollectNetWorkPictureStrategy{})
	h.AppendStrategy(&handler.CollectLocalPictureStrategy{})
	h.AppendStrategy(&InsertHexoHeaderStrategy{blogMarkdownRootDir})
	return h
}
