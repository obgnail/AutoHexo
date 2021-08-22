package AutoHexo

import (
	"github.com/fsnotify/fsnotify"
	"github.com/obgnail/AutoHexo/utils"
	"os"
)

//
//const MarkdownFileHeaderTemplate = "---\ntitle: %s\ndate: %s\ntags: %s\n---\n"
//
//func AddMdFileHeader(mdFilePath, title, date, tags string) error {
//	content := fmt.Sprintf(MarkdownFileHeaderTemplate, title, date, tags)
//	return InsertContentIntoFileHead(mdFilePath, []byte(content))
//}
//
//// TODO
//func RemoveMdFileHeader(mdFilePath string) error {
//	return nil
//}

func (ab *AutoBlogBuilder) BuildBlog(mdFilePath string) error {
	f, err := os.Open(mdFilePath)
	defer f.Close()
	if err != nil {
		return err
	}
	if err := utils.GetFileInfo(f); err != nil {
		return err
	}
	return nil
}

func (ab *AutoBlogBuilder) onFileCreate(events *utils.WatchEvents) {

}

func (ab *AutoBlogBuilder) onFileWrite(events *utils.WatchEvents) {

}

func (ab *AutoBlogBuilder) onFileRemove(events *utils.WatchEvents) {

}

func (ab *AutoBlogBuilder) onFileRename(events *utils.WatchEvents) {

}

func (ab *AutoBlogBuilder) onFileChmod(events *utils.WatchEvents) {

}

func (ab *AutoBlogBuilder) authHexo(ch chan *utils.EventMap) {
	for em := range ch {
		for typ, events := range em.m {
			switch typ {
			case fsnotify.Create.String():
				ab.onFileCreate(events)
			case fsnotify.Write.String():
				ab.onFileWrite(events)
			case fsnotify.Remove.String():
				ab.onFileRemove(events)
			case fsnotify.Rename.String():
				ab.onFileRename(events)
			case fsnotify.Chmod.String():
				ab.onFileChmod(events)
			}
		}
	}
}

func (ab *AutoBlogBuilder) Run() {
	watcher := utils.NewNotifyFile(ab.markdownDir, ab.waitingWindows, ab.authHexo)
	watcher.WatchDir()
	select {}
}
