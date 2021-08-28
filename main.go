package main

import (
	"github.com/obgnail/AutoHexo/notify_hexo"
	"time"
)

func main() {
	originMarkdownRootDir := `C:\Users\12516\Dropbox\root\md\Learning`
	blogMarkdownRootDir := `D:\app\blog\source\_posts`
	blogResourceRootDir := `D:\app\blog\source\images`
	hexoCmdPath := `C:\Users\12516\AppData\Roaming\npm\hexo`
	waitingWindows := 10 * time.Second

	autoHexo := notify_hexo.New(originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir, hexoCmdPath, waitingWindows)
	autoHexo.Run()
}
