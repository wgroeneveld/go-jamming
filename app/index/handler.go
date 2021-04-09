package index

import (
	"fmt"
	"net/http"

	"brainbaking.com/go-jamming/common"
)

func Handle(conf *common.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "This is a Jamstack microservice endpoint.\nWanna start jammin' too? Go to https://github.com/wgroeneveld/go-jamming !")
	}
}
