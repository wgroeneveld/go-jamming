
package app

import (
	"fmt"
    "strconv"
	"net/http"

	"github.com/gorilla/mux"

    "github.com/wgroeneveld/go-jamming/common"
)

type server struct {
	router *mux.Router
    conf *common.Config
}

// mimicing NotFound: https://golang.org/src/net/http/server.go?s=64787:64830#L2076
func unauthorized(w http.ResponseWriter, r *http.Request) { http.Error(w, "401 unauthorized", http.StatusUnauthorized) }	

func (s *server) authorizedOnly(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
	    vars := mux.Vars(r)
	    if vars["token"] != s.conf.Token {
        	unauthorized(w, r)
        	return
	    }
		h(w, r)
	}
}

func Start() {
    r := mux.NewRouter()
    server := &server{router: r, conf: common.Configure()}

    server.routes()
    http.Handle("/", r)

    fmt.Printf("Serving at port %d...\n", server.conf.Port)
    http.ListenAndServe(":" + strconv.Itoa(server.conf.Port), nil)
}
