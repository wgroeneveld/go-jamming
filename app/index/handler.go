
package index

import (
	"net/http"
	"fmt"

	"github.com/wgroeneveld/go-jamming/common"
)

func Handle(conf *common.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
	    w.WriteHeader(http.StatusOK)
	    fmt.Printf("testje")
    }
}

