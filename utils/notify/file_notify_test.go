package notify

import (
	"github.com/fsnotify/fsnotify"
	"testing"
	"time"
)

func TestFileNotify(t *testing.T) {
	waitingTime := 10 * time.Second
	path := "/Users/heyingliang/Dropbox/root/md"
	f := func(event fsnotify.Event) {
		t.Log(time.Now(), event)
	}
	watcher := New(path, waitingTime, f, nil)
	watcher.WatchDir()
	select {}
}
