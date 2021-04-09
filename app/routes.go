
package app

import (
	"brainbaking.com/go-jamming/app/index"
	"brainbaking.com/go-jamming/app/pingback"
	"brainbaking.com/go-jamming/app/webmention"
)

// stole ideas from https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
// not that contempt with passing conf, but can't create receivers on non-local types, and won't move specifics into package app
// https://blog.questionable.services/article/http-handler-error-handling-revisited/ is the better idea, but more work
func (s *server) routes() {
	cnf := s.conf

	s.router.HandleFunc("/", index.Handle(cnf)).Methods("GET")
	s.router.HandleFunc("/pingback", pingback.HandlePost(cnf)).Methods("POST")
	s.router.HandleFunc("/webmention", webmention.HandlePost(cnf)).Methods("POST")
	s.router.HandleFunc("/webmention/{domain}/{token}", s.authorizedOnly(webmention.HandleGet(cnf))).Methods("GET")
	s.router.HandleFunc("/webmention/{domain}/{token}", s.authorizedOnly(webmention.HandlePut(cnf))).Methods("PUT")
}

