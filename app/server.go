package app

import (
	"brainbaking.com/go-jamming/rest"
	"github.com/MagnusFrater/helmet"
	"net/http"
	"strconv"

	"brainbaking.com/go-jamming/common"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type server struct {
	router *mux.Router
	conf   *common.Config
}

func (s *server) authorizedOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["token"] != s.conf.Token || !s.conf.IsAnAllowedDomain(vars["domain"]) {
			rest.Unauthorized(w)
			return
		}
		h(w, r)
	}
}

func Start() {
	r := mux.NewRouter()
	config := common.Configure()
	config.SetupDataDirs()
	helmet := helmet.Default()

	server := &server{router: r, conf: config}

	server.routes()
	http.Handle("/", r)
	r.Use(LoggingMiddleware)
	r.Use(helmet.Secure)
	r.Use(NewRateLimiter(5, 10).Middleware)

	log.Info().Int("port", server.conf.Port).Msg("Serving...")
	http.ListenAndServe(":"+strconv.Itoa(server.conf.Port), nil)
}
