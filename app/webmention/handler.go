
package webmention

import (
	"github.com/wgroeneveld/go-jamming/app/mf"
	"github.com/wgroeneveld/go-jamming/app/webmention/receive"
	"net/http"
	"fmt"

	"github.com/wgroeneveld/go-jamming/common"
    "github.com/wgroeneveld/go-jamming/rest"
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
    httpClient := &rest.HttpClient{}

    return func(w http.ResponseWriter, r *http.Request) {
    	r.ParseForm()
    	if !validate(r, r.Header, conf) {
    		rest.BadRequest(w)
    		return
    	}
    	
    	target := r.FormValue("target")
    	if !isValidTargetUrl(target, httpClient) {
    		rest.BadRequest(w)
    		return
    	}

    	wm := mf.Mention{
            Source: r.FormValue("source"),
            Target: target,
        }
        recv := &receive.Receiver{
            RestClient: httpClient,
            Conf:       conf,
        }

        go recv.Receive(wm)
        rest.Accept(w)
    }
}

