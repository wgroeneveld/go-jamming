
package webmention

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wgroeneveld/go-jamming/app/mf"
	"github.com/wgroeneveld/go-jamming/app/webmention/receive"
	"github.com/wgroeneveld/go-jamming/app/webmention/send"
	"net/http"

	"github.com/wgroeneveld/go-jamming/common"
	"github.com/wgroeneveld/go-jamming/rest"
)

var httpClient = &rest.HttpClient{}

func HandleGet(conf *common.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    	fmt.Println("handling get")
    }
}

func HandlePut(conf *common.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		since := getSinceQueryParam(r)
		domain := mux.Vars(r)["domain"]

		snder := send.Sender{
			RestClient: httpClient,
			Conf:       conf,
		}
		go snder.Send(domain, since)
		rest.Accept(w)
	}
}

func getSinceQueryParam(r *http.Request) string {
	sinceParam, _ := r.URL.Query()["since"]
	since := ""
	if len(sinceParam) > 0 {
		since = sinceParam[0]
	}
	return since
}

func HandlePost(conf *common.Config) http.HandlerFunc {
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

