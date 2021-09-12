package auto_hexo

import (
	"runtime"
	"testing"
)

var (
	originMarkdownRootDir string
	blogMarkdownRootDir   string
	blogResourceRootDir   string
	hexoCmdPath           string
)

func TestAutoHexo(t *testing.T) {
	autoHexo := New(originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir,hexoCmdPath)
	if err := autoHexo.Run(originMarkdownRootDir); err != nil {
		t.Log("[WARN] hexo new err", err)
	}
}

func init() {
	if runtime.GOOS == "windows" {
		originMarkdownRootDir = `C:\Users\12516\Dropbox\root\md\Learning`
		blogMarkdownRootDir = `D:\app\blog\source\_posts`
		blogResourceRootDir = `D:\app\blog\source\images`
		hexoCmdPath = `C:\Users\12516\AppData\Roaming\npm\hexo`
	} else {
		originMarkdownRootDir = "/Users/heyingliang/myTemp/root/md2/Learning"
		blogMarkdownRootDir = "/Users/heyingliang/myTemp/blog/source/_posts"
		blogResourceRootDir = "/Users/heyingliang/myTemp/blog/source/images"
	}
}
