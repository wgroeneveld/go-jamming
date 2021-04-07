
package common

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "400 bad request", http.StatusBadRequest)
}

func Accept(w http.ResponseWriter) {
	w.WriteHeader(202)
	w.Write([]byte("Thanks, bro. Will send these webmentions soon, pinky swear!"))
}

// something like this? https://freshman.tech/snippets/go/http-response-to-string/
func Get(url string) (string, error) {
	resp, geterr := http.Get(url)
	if geterr != nil {
		return "", geterr
	}

    if resp.StatusCode < 200 || resp.StatusCode > 299 {
    	return "", fmt.Errorf("Status code for %s is not OK (%d)", url, resp.StatusCode)
    }

	defer resp.Body.Close()
	body, readerr := ioutil.ReadAll(resp.Body)
	if readerr != nil {
		return "", readerr
	}

	return string(body), nil
}
