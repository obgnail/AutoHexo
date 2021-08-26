package auto_hexo

import (
	"github.com/obgnail/AutoHexo/utils/notify"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/obgnail/AutoHexo/hexo_handler"
	"github.com/obgnail/AutoHexo/utils"
)

const (
	MarkdownFileSuffix     = ".md"
	LocalPictureUseAbsPath = true
)

type AutoHexo struct {
	watchMarkdownRootDir string        // 监控的md目录
	blogMarkdownRootDir  string        // 用于生成blog的目录
	blogResourceRootDir  string        // blog的资源目录
	waitingWindows       time.Duration // 等待窗口,合并若干时间内的消息
}

func New(watchMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir string, waitingWindows time.Duration) *AutoHexo {
	return &AutoHexo{
		watchMarkdownRootDir: watchMarkdownRootDir,
		blogMarkdownRootDir:  blogMarkdownRootDir,
		blogResourceRootDir:  blogResourceRootDir,
		waitingWindows:       waitingWindows,
	}
}

func (ah *AutoHexo) newBlog(mdFilesPath []string) error {
	handler := hexo_handler.New(
		ah.watchMarkdownRootDir, ah.blogResourceRootDir, ah.blogMarkdownRootDir,
		MarkdownFileSuffix, LocalPictureUseAbsPath, mdFilesPath,
	)
	handler.Run()
	return nil
}

// 删除dropbox在同步时产生的`~.md`临时文件
func (ah *AutoHexo) removeTempFiles(em *utils.EventMap) *utils.EventMap {
	return em.Filter(func(e fsnotify.Event) bool {
		filePath := notify.GetFilePathFromEvent(e)
		return !strings.HasSuffix(filePath, "~.md")
	})
}

func (ah *AutoHexo) reportNotifyFiles(em *utils.EventMap) {
	files := em.GetFileList()
	if len(files) != 0 {
		log.Println("--- changed file ---")
		for _, file := range files {
			log.Println("****** ", file)
		}
	}
}

func (ah *AutoHexo) RemoveFile(events *utils.WatchEvents) {
	for _, e := range events.Events {
		filePath := notify.GetFilePathFromEvent(e)
		blogFilePath := ah.GetBlogMarkdownPath(filePath)
		if err := os.Remove(blogFilePath); err != nil {
			log.Println("[WARN] file remove Error:", blogFilePath)
		} else {
			log.Print("[INFO] file remove:", blogFilePath)
		}
	}
}

// input:  /Users/heyingliang/Dropbox/root/md/Learning/Canal/网站__Cannal入门
// output: /Users/heyingliang/myTemp/root2/md4/Learning/Canal/网站__Cannal入门
func (ah *AutoHexo) GetBlogMarkdownPath(originMarkdownDir string) string {
	return strings.Replace(originMarkdownDir, ah.watchMarkdownRootDir, ah.blogMarkdownRootDir, 1)
}

func (ah *AutoHexo) onFileCreate(events *utils.WatchEvents) {
	var allCreateFilePath []string
	for _, e := range events.Events {
		filePath := notify.GetFilePathFromEvent(e)
		allCreateFilePath = append(allCreateFilePath, filePath)
	}
	if err := ah.newBlog(allCreateFilePath); err != nil {
		log.Println("[WARN] hexo new err", err)
	}
}

func (ah *AutoHexo) onFileWrite(events *utils.WatchEvents) {
	ah.onFileCreate(events)
}

func (ah *AutoHexo) onFileRemove(events *utils.WatchEvents) {
	ah.RemoveFile(events)
}

func (ah *AutoHexo) onFileRename(events *utils.WatchEvents) {
	ah.RemoveFile(events)
	ah.onFileCreate(events)
}

func (ah *AutoHexo) onFileChmod(events *utils.WatchEvents) {}

func (ah *AutoHexo) filterWatchDir(path string, info os.FileInfo) bool {
	basePath := filepath.Base(path)
	if basePath == "assets" || basePath == "__pycache__" {
		return false
	}
	return true
}

// 传过来一个chan,会根据不同的event type执行不同的hexo操作
func (ah *AutoHexo) authHexo(ch chan *utils.EventMap) {
	for em := range ch {
		em = ah.removeTempFiles(em)
		ah.reportNotifyFiles(em)
		for typ, events := range em.M {
			switch typ {
			case fsnotify.Create.String():
				ah.onFileCreate(events)
			case fsnotify.Write.String():
				ah.onFileWrite(events)
			case fsnotify.Remove.String():
				ah.onFileRemove(events)
			case fsnotify.Rename.String():
				ah.onFileRename(events)
			case fsnotify.Chmod.String():
				ah.onFileChmod(events)
			}
		}
	}
}

func (ah *AutoHexo) Run() {
	watcher := notify.New(ah.watchMarkdownRootDir, ah.waitingWindows, ah.authHexo, ah.filterWatchDir)
	watcher.WatchDir()
	select {}
}
