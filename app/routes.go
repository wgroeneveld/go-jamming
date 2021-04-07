
package app

import (
	"github.com/wgroeneveld/go-jamming/app/index"
	"github.com/wgroeneveld/go-jamming/app/webmention"
)

// stole bits from https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
// not that contempt with passing conf, but can't create receivers on non-local types, and won't move specifics into package app
func (s *server) routes() {
	s.router.HandleFunc("/", index.Handle(s.conf))
	s.router.HandleFunc("/webmention/{domain}/{token}", s.authorizedOnly(webmention.Handle(s.conf)))
}
