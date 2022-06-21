package rss

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	cnf = &common.Config{
		BaseURL:                  "http://localhost:1337/",
		Port:                     1337,
		Token:                    "miauwkes",
		AllowedWebmentionSources: []string{"brainbaking.com"},
		Blacklist:                []string{},
		Whitelist:                []string{"brainbaking.com"},
	}
	repo db.MentionRepo
)

func init() {
	repo = db.NewMentionRepo(cnf)
}
func TestHandleGet(t *testing.T) {
	wmInMod := mf.Mention{
		Source: "https://infos.by/markdown-v-nauke/",
		Target: "https://brainbaking.com/post/2021/02/writing-academic-papers-in-markdown/",
	}
	wmApproved := mf.Mention{
		Source: "https://brainbaking.com/post/2022/04/equality-in-game-credits/",
		Target: "https://brainbaking.com/",
	}

	repo.InModeration(wmInMod, &mf.IndiewebData{
		Source: wmInMod.Source,
		Target: wmInMod.Target,
		Name:   "inmod1",
	})
	repo.Save(wmApproved, &mf.IndiewebData{
		Source: wmApproved.Source,
		Target: wmApproved.Target,
		Name:   "approved1",
	})
	r := mux.NewRouter()
	r.HandleFunc("/feed/{domain}/{token}", HandleGet(cnf, repo)).Methods("GET")
	ts := httptest.NewServer(r)

	t.Cleanup(func() {
		os.Remove("config.json")
		ts.Close()
		db.Purge()
	})

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/feed/%s/%s", ts.URL, cnf.AllowedWebmentionSources[0], cnf.Token), nil)

	resp, err := client.Do(req)
	assert.NoError(t, err)

	contentBytes, _ := ioutil.ReadAll(resp.Body)
	content := string(contentBytes)
	defer resp.Body.Close()

	assert.Contains(t, content, "<description>Go-Jamming @ brainbaking.com</description>")
	assert.Contains(t, content, "<title>To Moderate: inmod1 ()</title>")
	assert.Contains(t, content, "<title>approved1 ()</title>")
}
