package strategy

import (
	"fmt"
	"path/filepath"

	"github.com/obgnail/MarkdownResouceCollecter/handler"
)

// input: /Users/heyingliang/myTemp/blog/source/images/sync-trade-off.svg
// output: /images/sync-trade-off.svg
type FixHexoPicturePathStrategy struct{}

func (s *FixHexoPicturePathStrategy) BeforeRewrite(h *handler.BaseHandler) error {
	for _, file := range h.Files {
		for _, pic := range file.Pictures {
			// 网络图片和不存在的图片保持不变
			if pic.FromNet || !pic.IsExist {
				pic.NewPath = pic.OriginPath
				pic.NewMatch = pic.OriginMatch
				continue
			}
			fileBase := filepath.Base(pic.NewPath)
			imageBase := filepath.Base(filepath.Dir(pic.NewPath))
			// hexo强制使用![](/images/image.jpg)来访问资源
			pic.NewPath = fmt.Sprintf("/%s/%s",imageBase, fileBase)
			pic.NewMatch = fmt.Sprintf("![%s](%s)", pic.ShowName, pic.NewPath)
		}
	}
	return nil
}

func (s *FixHexoPicturePathStrategy) AfterRewrite(h *handler.BaseHandler) error { return nil }
