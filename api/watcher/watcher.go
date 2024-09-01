package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/xiaoxuan6/deeplx"
	"github.com/xiaoxuan6/deeplx/api/log"
)

func WatchBlackList() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("%s", err.Error())
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Errorf("watcher error：%s", err.Error())
					return
				}
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if (event.Op & fsnotify.Write) != 0 {
					deeplx.LoadBlack(true)
				}
			}
		}
	}()

	_ = watcher.Add("blacklist.txt")
	<-make(chan struct{})
}
