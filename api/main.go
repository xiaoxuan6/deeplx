package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

	r := mux.NewRouter()
	r.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFound)

	route.Register(r)
	r.HandleFunc("/", handlers.Index)

	_ = http.ListenAndServe(":8311", r)
}
