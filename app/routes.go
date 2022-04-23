package app

import (
	"brainbaking.com/go-jamming/app/admin"
	"brainbaking.com/go-jamming/app/index"
	"brainbaking.com/go-jamming/app/pictures"
	"brainbaking.com/go-jamming/app/pingback"
	"brainbaking.com/go-jamming/app/webmention"
)

// stole ideas from https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
// not that contempt with passing conf, but can't create receivers on non-local types, and won't move specifics into package app
// https://blog.questionable.services/article/http-handler-error-handling-revisited/ is the better idea, but more work
func (s *server) routes() {
	c := s.conf
	db := s.repo

	s.router.HandleFunc("/", index.Handle(c)).Methods("GET")
	s.router.HandleFunc("/pictures/{picture}", pictures.Handle(db)).Methods("GET")
	s.router.HandleFunc("/pingback", pingback.HandlePost(c, db)).Methods("POST")
	s.router.HandleFunc("/webmention", webmention.HandlePost(c, db)).Methods("POST")

	s.router.HandleFunc("/webmention/{domain}/{token}", s.domainAndTokenOnly(webmention.HandleGet(db))).Methods("GET")
	s.router.HandleFunc("/webmention/{domain}/{token}", s.domainAndTokenOnly(webmention.HandlePut(c, db))).Methods("PUT")
	s.router.HandleFunc("/webmention/{domain}/{token}", s.domainAndTokenOnly(webmention.HandleDelete(db))).Methods("DELETE")

	s.router.HandleFunc("/admin/{domain}/{token}", s.domainAndTokenOnly(admin.HandleGet(db))).Methods("GET")
	s.router.HandleFunc("/admin/approve/{token}/{key}", s.authorizedOnly(admin.HandleApprove(c, db))).Methods("GET")
	s.router.HandleFunc("/admin/reject/{token}/{key}", s.authorizedOnly(admin.HandleReject(c, db))).Methods("GET")
}
