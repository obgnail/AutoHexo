package notify_hexo

import (
	"github.com/obgnail/AutoHexo/auto_hexo"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/obgnail/AutoHexo/utils/notify"
)

const (
	unixTempFileSuffix    = "~.md"
	windowsTempFilePrefix = ".~"
	markdownAssetsDir     = "assets"
	pythonCacheDir        = "__pycache__"
	macDSFile             = ".DS_Store"
)

type NotifyHexo struct {
	*auto_hexo.AutoHexo
	waitingWindows time.Duration // 等待窗口,合并若干时间内的消息
}

func New(
	originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir, hexoCmdPath string,
	waitingWindows time.Duration,
) *NotifyHexo {
	ah := auto_hexo.New(originMarkdownRootDir, blogMarkdownRootDir, blogResourceRootDir, hexoCmdPath)
	return &NotifyHexo{AutoHexo: ah, waitingWindows: waitingWindows}
}

func (nh *NotifyHexo) newBlog(changedFilePath string) error {
	return nh.CreateBlog(changedFilePath)
}

// 删除同步时产生的`filename~.md`(unix)和`.~filename.md`(windows)临时文件
func (nh *NotifyHexo) isTempFile(event fsnotify.Event) (ret bool) {
	if runtime.GOOS == "windows" {
		filePath := notify.GetFilePathFromEvent(event)
		filePath = filepath.Base(filePath)
		ret = strings.HasPrefix(filePath, windowsTempFilePrefix)
	} else {
		filePath := notify.GetFilePathFromEvent(event)
		ret = strings.HasSuffix(filePath, unixTempFileSuffix)
	}
	return
}

func (nh *NotifyHexo) removeBlogFile(event fsnotify.Event) {
	removeFilePath := notify.GetFilePathFromEvent(event)
	blogFilePath := strings.Replace(removeFilePath, nh.OriginMarkdownRootDir, nh.BlogMarkdownRootDir, 1)
	if err := nh.DeleteBlog(blogFilePath); err != nil {
		log.Println("[WARN] file remove Error:", blogFilePath)
	}
}

func (nh *NotifyHexo) onFileCreate(event fsnotify.Event) {
	changedFilePath := notify.GetFilePathFromEvent(event)
	if err := nh.newBlog(changedFilePath); err != nil {
		log.Println("[WARN] hexo new err", err)
	}
}

func (nh *NotifyHexo) onFileWrite(event fsnotify.Event) {
	nh.onFileCreate(event)
}

func (nh *NotifyHexo) onFileRemove(event fsnotify.Event) {
	nh.removeBlogFile(event)
}

func (nh *NotifyHexo) onFileRename(event fsnotify.Event) {
	nh.removeBlogFile(event)
	nh.onFileCreate(event)
}

func (nh *NotifyHexo) onFileChmod(event fsnotify.Event) {}

func (nh *NotifyHexo) filterWatchDir(path string, info os.FileInfo) bool {
	basePath := filepath.Base(path)
	return !(basePath == markdownAssetsDir || basePath == pythonCacheDir || basePath == macDSFile)
}

// AuthHexo 根据不同的event type执行不同的hexo操作
func (nh *NotifyHexo) AuthHexo(event fsnotify.Event) {
	if nh.isTempFile(event) {
		return
	}
	if notify.IsEventBelongToType(event, fsnotify.Create) {
		nh.onFileCreate(event)
	}
	if notify.IsEventBelongToType(event, fsnotify.Write) {
		nh.onFileWrite(event)
	}
	if notify.IsEventBelongToType(event, fsnotify.Remove) {
		nh.onFileRemove(event)
	}
	if notify.IsEventBelongToType(event, fsnotify.Rename) {
		nh.onFileRename(event)
	}
	if notify.IsEventBelongToType(event, fsnotify.Chmod) {
		nh.onFileChmod(event)
	}
}

func (nh *NotifyHexo) Run() {
	watcher := notify.New(nh.OriginMarkdownRootDir, nh.waitingWindows, nh.AuthHexo, nh.filterWatchDir)
	watcher.WatchDir()
	select {}
}
