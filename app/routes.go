
package app

import (
	"github.com/wgroeneveld/go-jamming/app/index"
)

func (s *server) routes() {
	s.router.HandleFunc("/", index.HandleIndex)
}
