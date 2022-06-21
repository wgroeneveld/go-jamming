package app

import (
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
	"github.com/goddtriffin/helmet"
	"net/http"
	"strconv"
	"strings"

	"brainbaking.com/go-jamming/common"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type server struct {
	router *mux.Router
	conf   *common.Config
	repo   db.MentionRepo
}

func (s *server) domainAndTokenOnly(h http.HandlerFunc) http.HandlerFunc {
	return s.domainOnly(s.tokenOnly(h))
}

func (s *server) domainOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if !s.conf.IsAnAllowedDomain(vars["domain"]) {
			rest.Unauthorized(w)
			return
		}
		h(w, r)
	}
}

func (s *server) tokenOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["token"] != s.conf.Token {
			rest.Unauthorized(w)
			return
		}
		h(w, r)
	}
}

func ipFrom(r *http.Request) string {
	realIp := r.Header.Get("X-Real-IP")
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if realIp != "" { // in case of proxy. is IP itself
		return realIp
	}
	if forwardedFor != "" { // in case of proxy. Could be: clientip, proxy1, proxy2, ...
		return strings.Split(forwardedFor, ",")[0]
	}
	if strings.Contains(r.RemoteAddr, ":") { // in case of 127.0.0.1:12345
		return strings.Split(r.RemoteAddr, ":")[0]
	}
	return r.RemoteAddr
}

func Start() {
	r := mux.NewRouter()
	config := common.Configure()
	repo := db.NewMentionRepo(config)
	helmet := helmet.Default()

	server := &server{router: r, conf: config, repo: repo}

	server.routes()
	http.Handle("/", r)
	r.Use(LoggingMiddleware)
	r.Use(helmet.Secure)
	r.Use(NewRateLimiter(5, 10).Middleware)

	log.Info().Int("port", server.conf.Port).Msg("Serving...")
	log.Fatal().Err(http.ListenAndServe(":"+strconv.Itoa(server.conf.Port), nil))
}
