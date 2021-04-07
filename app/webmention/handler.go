
package webmention

import (
	"net/http"
	"fmt"

	"github.com/wgroeneveld/go-jamming/common"
)

func HandleGet(conf *common.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    	fmt.Println("handling get")
    }
}

func HandlePut(conf *common.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    	fmt.Println("handling put")
    }
}

func HandlePost(conf *common.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    }
}

