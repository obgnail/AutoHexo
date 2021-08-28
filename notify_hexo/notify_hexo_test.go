package notify_hexo

import (
	"runtime"
	"testing"
	"time"
)

var (
	originMarkdownRootDir string
	blogMarkdownRootDir   string
	blogResourceRootDir   string
	hexoCmdPath           string
	waitingWindows        time.Duration
)

func TestNotifyHexo(t *testing.T) {
	notifyHexo := New(originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir, hexoCmdPath, waitingWindows)
	notifyHexo.Run()
	select {}
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
	waitingWindows = 10 * time.Second
}
