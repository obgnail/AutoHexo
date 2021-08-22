package AutoHexo

import "time"

type AutoBlogBuilder struct {
	markdownDir string

	// 等待窗口,合并若干时间内的消息
	waitingWindows time.Duration
}
