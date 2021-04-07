
package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
}

func Start() {
    r := mux.NewRouter()
    server := &server{router: r}

    server.routes()
    http.Handle("/", r)

    fmt.Printf("Serving at port 1337...\n")
    http.ListenAndServe(":1337", nil)
}
