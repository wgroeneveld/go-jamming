package webmention

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/webmention/recv"
	"brainbaking.com/go-jamming/app/webmention/send"
	"brainbaking.com/go-jamming/db"
	"github.com/gorilla/mux"
	"net/http"

	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/rest"
)

var (
	httpClient = &rest.HttpClient{}
)

func HandleGet(repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := mux.Vars(r)["domain"]
		rest.Json(w, repo.GetAll(domain))
	}
}

func HandlePut(conf *common.Config, repo db.MentionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := mux.Vars(r)["domain"]
		source := sourceQueryParam(r)

		snder := &send.Sender{
			RestClient: httpClient,
			Conf:       conf,
			Repo:       repo,
		}

		if source != "" {
			go snder.SendSingle(domain, source)
		} else {
			go snder.Send(domain)
		}

		rest.Accept(w)
	}
}

func sourceQueryParam(r *http.Request) string {
	sourceParam := r.URL.Query()["source"]
	if len(sourceParam) > 0 {
		return sourceParam[0]
	}
	return ""
}

func HandlePost(conf *common.Config, repo db.MentionRepo) http.HandlerFunc {
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
			Repo:       repo,
		}

		go recv.Receive(wm)
		rest.Accept(w)
	}
}
