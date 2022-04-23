package app

import (
	"brainbaking.com/go-jamming/common"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var conf = &common.Config{
	Token:                    "boemsjakkalakka",
	AllowedWebmentionSources: []string{"http://ewelja.be"},
}

func TestAuthorizedOnlyUnauthorizedWithWrongToken(t *testing.T) {
	srv := &server{
		conf: conf,
	}

	passed := false
	handler := srv.authorizedOnly(func(writer http.ResponseWriter, request *http.Request) {
		passed = true
	})
	r, _ := http.NewRequest("PUT", "/whatever", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{
		"token":  "invalid",
		"domain": conf.AllowedWebmentionSources[0],
	})

	handler(w, r)
	assert.False(t, passed, "should not have called unauthorized func")
}

func TestDomainOnlyWithWrongDomain(t *testing.T) {
	srv := &server{
		conf: conf,
	}

	passed := false
	handler := srv.domainOnly(func(writer http.ResponseWriter, request *http.Request) {
		passed = true
	})
	r, _ := http.NewRequest("PUT", "/whatever", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{
		"token":  conf.Token,
		"domain": "https://sexymoddafokkas.be",
	})

	handler(w, r)
	assert.False(t, passed, "should not have called unauthorized func")
}

func TestAuthorizedOnlyOkIfTokenAndDomainMatch(t *testing.T) {
	srv := &server{
		conf: conf,
	}

	passed := false
	handler := srv.authorizedOnly(func(writer http.ResponseWriter, request *http.Request) {
		passed = true
	})
	r, _ := http.NewRequest("PUT", "/whatever", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{
		"token":  conf.Token,
		"domain": conf.AllowedWebmentionSources[0],
	})

	handler(w, r)
	assert.True(t, passed, "should have passed authentication!")
}
