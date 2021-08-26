package main

import (
	"github.com/obgnail/AutoHexo/auto_hexo"
	"time"
)

func main() {
	originMarkdownRootDir := "/Users/heyingliang/myTemp/root/md2/Learning"
	blogMarkdownRootDir := "/Users/heyingliang/myTemp/blog/source/_posts"
	blogResourceRootDir := "/Users/heyingliang/myTemp/blog/source/images"
	waitingWindows := 10 * time.Second

	autoHexo := auto_hexo.New(originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir, waitingWindows)
	autoHexo.Run()
}
