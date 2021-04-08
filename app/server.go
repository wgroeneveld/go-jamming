
package app

import (
    "strconv"
	"net/http"

    "github.com/wgroeneveld/go-jamming/common"

    "github.com/gorilla/mux"
    "github.com/rs/zerolog/log"
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
    config := common.Configure()
    config.SetupDataDirs()
    server := &server{router: r, conf: config}

    server.routes()
    http.Handle("/", r)
    r.Use(loggingMiddleware)

    log.Info().Int("port", server.conf.Port).Msg("Serving...")
    http.ListenAndServe(":" + strconv.Itoa(server.conf.Port), nil)
}
