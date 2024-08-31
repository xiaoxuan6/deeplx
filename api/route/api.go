package route

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/xiaoxuan6/deeplx/api/handlers"
	"os"
	"strings"
)

func Register(r *mux.Router) {
	path := "translate"
	if routerPath := os.Getenv("ROUTER_PATH"); len(routerPath) > 0 {
		path = routerPath
	}

	path = fmt.Sprintf("/%s", strings.Trim(path, "/"))
	r.HandleFunc(path, handlers.Translate).Methods("POST")
}
