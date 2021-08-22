package utils

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var eventMap *EventMap

type WatchEvents struct {
	events []fsnotify.Event
}

type EventMap struct {
	*sync.Mutex
	m map[string]*WatchEvents
}

func NewEventMap() *EventMap {
	return &EventMap{
		Mutex: &sync.Mutex{},
		m:     make(map[string]*WatchEvents, 5),
	}
}

func (em *EventMap) Append(event fsnotify.Event) {
	typeArray := []fsnotify.Op{
		fsnotify.Create, fsnotify.Write, fsnotify.Remove, fsnotify.Rename, fsnotify.Chmod,
	}
	em.Lock()
	defer em.Unlock()
	// 一个event可能属于多种type
	for _, typ := range typeArray {
		if isEventBelongToType(event, typ) {
			name := typ.String()
			if _, ok := em.m[name]; !ok {
				em.m[name] = new(WatchEvents)
			}
			em.m[name].events = append(em.m[name].events, event)
		}
	}
}

func (em *EventMap) GetFileList() []string {
	var fileList []string
	files := make(map[string]struct{}, 0)
	em.Lock()
	defer em.Unlock()
	for _, events := range em.m {
		for _, event := range events.events {
			if _, exist := files[event.Name]; !exist {
				files[event.Name] = struct{}{}
				fileList = append(fileList, event.Name)
			}
		}
	}
	return fileList
}

func (em *EventMap) Clear() {
	em.Lock()
	defer em.Unlock()
	em.m = map[string]*WatchEvents{}
}

func (em *EventMap) Copy() *EventMap {
	ret := NewEventMap()
	for k, v := range em.m {
		ret.m[k] = v
	}
	return ret
}

type FileWatcher struct {
	watch           *fsnotify.Watcher
	eventInputChan  chan fsnotify.Event
	eventOutputChan chan *EventMap // 消息整合

	path            string
	waitingPushTime time.Duration
	callbackFunc    func(ch chan *EventMap)
}

func NewNotifyFile(path string, waitingPushTime time.Duration, callbackFunc func(ch chan *EventMap)) *FileWatcher {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	fw := &FileWatcher{
		watch:           watch,
		eventInputChan:  make(chan fsnotify.Event, 100),
		eventOutputChan: make(chan *EventMap, 1), // 在waitingPushTime内只发一条event

		path:            path,
		waitingPushTime: waitingPushTime,
		callbackFunc:    callbackFunc,
	}
	go fw.callback()
	return fw
}

func (fw *FileWatcher) callback() {
	fw.callbackFunc(fw.eventOutputChan)
}

func (fw *FileWatcher) pushEventMapToOutputChan() {
	ticker := time.NewTicker(fw.waitingPushTime)
	go func() {
		for range ticker.C {
			e := eventMap.Copy()
			fw.eventOutputChan <- e
			eventMap.Clear()
		}
	}()
}

func (fw *FileWatcher) collectEventToMap() {
	for {
		select {
		case event := <-fw.eventInputChan:
			eventMap.Append(event)
		}
	}
}

func (fw *FileWatcher) WatchDir() {
	filepath.Walk(fw.path, func(path string, info os.FileInfo, err error) error {
		// 因为目录下文件也在监控范围内,不需要加文件
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = fw.watch.Add(path)
			if err != nil {
				return err
			}
			log.Println("[INFO] add watch: ", path)
		}
		return nil
	})
	go fw.pushEventMapToOutputChan()
	go fw.collectEventToMap()
	go fw.watchEvent()
}

func (fw *FileWatcher) watchEvent() {
	for {
		select {
		case ev := <-fw.watch.Events:
			fw.eventInputChan <- ev
			if isEventBelongToType(ev, fsnotify.Create) {
				fw.onCreate(ev)
			}
			if isEventBelongToType(ev, fsnotify.Write) {
				fw.onWrite(ev)
			}
			if isEventBelongToType(ev, fsnotify.Remove) {
				fw.onRemove(ev)
			}
			if isEventBelongToType(ev, fsnotify.Rename) {
				fw.onRename(ev)
			}
			if isEventBelongToType(ev, fsnotify.Chmod) {
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

// JudgeEventType 判断event是否属于typ类型
func isEventBelongToType(event fsnotify.Event, typ fsnotify.Op) bool {
	return event.Op&typ == typ
}

func init() {
	eventMap = NewEventMap()
}
