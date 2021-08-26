package notify

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

const inputChanSize = 100

// WatchDirFilterFunc return false if you want filter this dir
type WatchDirFilterFunc func(path string, info os.FileInfo) bool
type CallbackFunc func(event fsnotify.Event)

type FileWatcher struct {
	path            string
	waitingPushTime time.Duration
	callbackFunc    CallbackFunc
	dirFilterFunc   WatchDirFilterFunc

	watch         *fsnotify.Watcher
	tickerChannel *TickerChannel
}

func New(
	path string, waitingPushTime time.Duration,
	callbackFunc CallbackFunc, dirFilterFunc WatchDirFilterFunc,
) *FileWatcher {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	tc := NewTickerChannel(inputChanSize, waitingPushTime)
	fw := &FileWatcher{
		watch:         watch,
		tickerChannel: tc,
		path:          path,
		dirFilterFunc: dirFilterFunc,
		callbackFunc:  callbackFunc,
	}
	go fw.handler()
	return fw
}

func (fw *FileWatcher) handler() {
	go func() {
		fw.tickerChannel.Range(func(event interface{}) {
			e := event.(fsnotify.Event)
			fw.callbackFunc(e)
		})
	}()
}

func (fw *FileWatcher) WatchDir() {
	filepath.Walk(fw.path, func(path string, info os.FileInfo, err error) error {
		if fw.dirFilterFunc != nil && !fw.dirFilterFunc(path, info) {
			return nil
		}
		// 因为目录下文件也在监控范围内,不需要加文件
		if info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = fw.watch.Add(absPath)
			if err != nil {
				return err
			}
			log.Println("[INFO] add watch: ", absPath)
		}
		return nil
	})
	go fw.watchEvent()
}

func (fw *FileWatcher) watchEvent() {
	for {
		select {
		case ev := <-fw.watch.Events:
			fw.tickerChannel.Store(ev)
			if IsEventBelongToType(ev, fsnotify.Create) {
				fw.onCreate(ev)
			}
			if IsEventBelongToType(ev, fsnotify.Write) {
				fw.onWrite(ev)
			}
			if IsEventBelongToType(ev, fsnotify.Remove) {
				fw.onRemove(ev)
			}
			if IsEventBelongToType(ev, fsnotify.Rename) {
				fw.onRename(ev)
			}
			if IsEventBelongToType(ev, fsnotify.Chmod) {
				fw.onChmod(ev)
			}
		case err := <-fw.watch.Errors:
			fw.onError(err)
			return
		}
	}
}

func (fw *FileWatcher) onCreate(event fsnotify.Event) {
	log.Println("[INFO] create: ", event.Name)
	// 获取新创建文件的信息,如果是目录,则加入监控中
	file, err := os.Stat(event.Name)
	if err == nil && file.IsDir() {
		if err := fw.watch.Add(event.Name); err != nil {
			log.Println("[WARN] add watch failed: ", event.Name)
			return
		}
		log.Println("[INFO] add watch: ", event.Name)
	}
}

func (fw *FileWatcher) onWrite(event fsnotify.Event) {
	log.Println("[INFO] write: ", event.Name)
}

func (fw *FileWatcher) onRemove(event fsnotify.Event) {
	log.Println("[INFO] remove: ", event.Name)
	// 如果删除文件是目录,则移除监控
	fi, err := os.Stat(event.Name)
	if err == nil && fi.IsDir() {
		if err := fw.watch.Remove(event.Name); err != nil {
			log.Println("[WARN] delete watch failed: ", event.Name)
			return
		}
		log.Println("[INFO] delete watch: ", event.Name)
	}
}

func (fw *FileWatcher) onRename(event fsnotify.Event) {
	log.Println("[INFO] rename: ", event.Name)
	// 如果重命名文件是目录,则移除监控,注意这里无法使用os.Stat来判断是否是目录了
	// 因为重命名后,go已经无法找到原文件来获取信息了,所以简单粗爆直接remove
	if err := fw.watch.Remove(event.Name); err != nil {
		log.Println("[WARN] delete watch failed: ", event.Name)
	}
}

func (fw *FileWatcher) onChmod(event fsnotify.Event) {
	log.Println("[INFO] chmod: ", event.Name)
}

func (fw *FileWatcher) onError(err error) {
	log.Println("[ERROR] watcher on error: ", err)
}

// IsEventBelongToType 判断event是否属于typ类型
func IsEventBelongToType(event fsnotify.Event, typ fsnotify.Op) bool {
	return event.Op&typ == typ
}

func GetFilePathFromEvent(event fsnotify.Event) string {
	return event.Name
}
