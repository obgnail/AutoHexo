package utils

import (
	"testing"
	"time"
)

func TestFileNotify(t *testing.T) {
	waitingTime := 10 * time.Second
	path := "/Users/heyingliang/Dropbox/root/md"
	f := func(ch chan *EventMap) {
		for em := range ch {
			files := em.GetFileList()
			t.Log(time.Now(), files)
		}
	}
	watcher := NewNotifyFile(path, waitingTime, f)
	watcher.WatchDir()
	select {}
}
