package webmention

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/webmention/load"
	"brainbaking.com/go-jamming/app/webmention/recv"
	"brainbaking.com/go-jamming/app/webmention/send"
	"github.com/gorilla/mux"
	"net/http"

	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/rest"
)

var httpClient = &rest.HttpClient{}

func HandleGet(conf *common.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := mux.Vars(r)["domain"]
		result := load.FromDisk(domain, conf.DataPath)

		rest.Json(w, result)
	}
}

func HandlePut(conf *common.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		since := getSinceQueryParam(r)
		domain := mux.Vars(r)["domain"]

		snder := &send.Sender{
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
		recv := &recv.Receiver{
			RestClient: httpClient,
			Conf:       conf,
		}

		go recv.Receive(wm)
		rest.Accept(w)
	}
}
