package auto_hexo

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/obgnail/AutoHexo/hexo_handler"
	"github.com/obgnail/AutoHexo/utils/notify"
)

const (
	MarkdownFileSuffix     = ".md"
	LocalPictureUseAbsPath = true

	TempFileSuffix    = "~.md"
	MarkdownAssetsDir = "assets"
	PythonCacheDir    = "__pycache__"
)

type AutoHexo struct {
	originMarkdownRootDir string        // 监控的md目录
	blogMarkdownRootDir   string        // 用于生成blog的目录
	blogResourceRootDir   string        // blog的资源目录
	waitingWindows        time.Duration // 等待窗口,合并若干时间内的消息
}

func New(
	originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir string,
	waitingWindows time.Duration,
) *AutoHexo {
	return &AutoHexo{
		originMarkdownRootDir: originMarkdownRootDir,
		blogMarkdownRootDir:   blogMarkdownRootDir,
		blogResourceRootDir:   blogResourceRootDir,
		waitingWindows:        waitingWindows,
	}
}

func (ah *AutoHexo) newBlog(changedFilePath string) error {
	h := hexo_handler.New(
		ah.originMarkdownRootDir, ah.blogResourceRootDir, ah.blogMarkdownRootDir,
		MarkdownFileSuffix, LocalPictureUseAbsPath, changedFilePath,
	)
	return h.Run()
}

// 删除同步时产生的`~.md`临时文件
func (ah *AutoHexo) isTempFile(event fsnotify.Event) bool {
	filePath := notify.GetFilePathFromEvent(event)
	return strings.HasSuffix(filePath, TempFileSuffix)
}

func (ah *AutoHexo) removeBlogFile(event fsnotify.Event) {
	removeFilePath := notify.GetFilePathFromEvent(event)
	blogFilePath := strings.Replace(removeFilePath, ah.originMarkdownRootDir, ah.blogMarkdownRootDir, 1)
	if err := os.Remove(blogFilePath); err != nil {
		log.Println("[WARN] file remove Error:", blogFilePath)
	} else {
		log.Print("[INFO] file remove:", blogFilePath)
	}
}

func (ah *AutoHexo) onFileCreate(event fsnotify.Event) {
	changedFilePath := notify.GetFilePathFromEvent(event)
	if err := ah.newBlog(changedFilePath); err != nil {
		log.Println("[WARN] hexo new err", err)
	}
}

func (ah *AutoHexo) onFileWrite(event fsnotify.Event) {
	ah.onFileCreate(event)
}

func (ah *AutoHexo) onFileRemove(event fsnotify.Event) {
	ah.removeBlogFile(event)
}

func (ah *AutoHexo) onFileRename(event fsnotify.Event) {
	ah.removeBlogFile(event)
	ah.onFileCreate(event)
}

func (ah *AutoHexo) onFileChmod(event fsnotify.Event) {}

func (ah *AutoHexo) filterWatchDir(path string, info os.FileInfo) bool {
	basePath := filepath.Base(path)
	return !(basePath == MarkdownAssetsDir || basePath == PythonCacheDir)
}

// 根据不同的event type执行不同的hexo操作
func (ah *AutoHexo) authHexo(event fsnotify.Event) {
	if ah.isTempFile(event) {
		return
	}
	if notify.IsEventBelongToType(event, fsnotify.Create) {
		ah.onFileCreate(event)
	}
	if notify.IsEventBelongToType(event, fsnotify.Write) {
		ah.onFileWrite(event)
	}
	if notify.IsEventBelongToType(event, fsnotify.Remove) {
		ah.onFileRemove(event)
	}
	if notify.IsEventBelongToType(event, fsnotify.Rename) {
		ah.onFileRename(event)
	}
	if notify.IsEventBelongToType(event, fsnotify.Chmod) {
		ah.onFileChmod(event)
	}
}

func (ah *AutoHexo) Run() {
	watcher := notify.New(ah.originMarkdownRootDir, ah.waitingWindows, ah.authHexo, ah.filterWatchDir)
	watcher.WatchDir()
	select {}
}
