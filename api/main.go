package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/xiaoxuan6/deeplx"
	"github.com/xiaoxuan6/deeplx/api/handlers"
	"github.com/xiaoxuan6/deeplx/api/log"
	"github.com/xiaoxuan6/deeplx/api/route"
	"net/http"
)

func init() {
	log.InitLog()
}

func main() {
	_ = godotenv.Load()
	deeplx.LoadBlack(false)

	go watchBlackList()

	r := mux.NewRouter()
	r.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFound)

	route.Register(r)
	r.HandleFunc("/", handlers.Index)

	_ = http.ListenAndServe(":8311", r)
}

func watchBlackList() {
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
