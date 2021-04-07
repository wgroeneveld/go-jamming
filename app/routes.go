
package app

import (
	"github.com/wgroeneveld/go-jamming/app/index"
	"github.com/wgroeneveld/go-jamming/app/webmention"
	"github.com/wgroeneveld/go-jamming/app/pingback"
)

// stole ideas from https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
// not that contempt with passing conf, but can't create receivers on non-local types, and won't move specifics into package app
func (s *server) routes() {
	s.router.HandleFunc("/", index.Handle(s.conf)).Methods("GET")
	s.router.HandleFunc("/pingback", pingback.Handle(s.conf)).Methods("POST")
	s.router.HandleFunc("/webmention", webmention.HandlePost(s.conf)).Methods("POST")
	s.router.HandleFunc("/webmention/{domain}/{token}", s.authorizedOnly(webmention.HandleGet(s.conf))).Methods("GET")
	s.router.HandleFunc("/webmention/{domain}/{token}", s.authorizedOnly(webmention.HandlePut(s.conf))).Methods("PUT")
}
