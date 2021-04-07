
package pingback

import (
	"net/http"

	"github.com/wgroeneveld/go-jamming/common"
)

func Handle(conf *common.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    }
}

