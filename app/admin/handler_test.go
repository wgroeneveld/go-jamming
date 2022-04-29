package admin

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
	wm := mf.Mention{
		Source: "https://infos.by/markdown-v-nauke/",
		Target: "https://brainbaking.com/post/2021/02/writing-academic-papers-in-markdown/",
	}

	repo.InModeration(wm, &mf.IndiewebData{
		Source: wm.Source,
		Target: wm.Target,
		Name:   "mytest",
	})
	r := mux.NewRouter()
	r.HandleFunc("/admin/{token}", HandleGet(cnf, repo)).Methods("GET")
	ts := httptest.NewServer(r)

	t.Cleanup(func() {
		os.Remove("config.json")
		ts.Close()
		db.Purge()
	})

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/admin/%s", ts.URL, cnf.Token), nil)

	resp, err := client.Do(req)
	assert.NoError(t, err)

	contentBytes, _ := ioutil.ReadAll(resp.Body)
	content := string(contentBytes)
	defer resp.Body.Close()

	assert.Contains(t, content, "admin dashboard")
	assert.Contains(t, content, wm.Source)
	assert.Contains(t, content, wm.Target)
}

func TestHandleReject(t *testing.T) {
	wm := mf.Mention{
		Source: "https://infos.by/markdown-v-nauke/",
		Target: "https://brainbaking.com/post/2021/02/writing-academic-papers-in-markdown/",
	}

	key, _ := repo.InModeration(wm, &mf.IndiewebData{
		Source: wm.Source,
		Target: wm.Target,
		Name:   "mytest",
	})
	assert.NotEmpty(t, repo.GetAllToModerate("brainbaking.com").Data)

	r := mux.NewRouter()
	r.HandleFunc("/admin/reject/{token}/{key}", HandleReject(cnf, repo)).Methods("GET")
	ts := httptest.NewServer(r)

	t.Cleanup(func() {
		os.Remove("config.json")
		ts.Close()
		db.Purge()
	})

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/admin/reject/%s/%s", ts.URL, cnf.Token, key), nil)

	_, err = client.Do(req)
	assert.NoError(t, err)

	assert.Empty(t, repo.GetAllToModerate("brainbaking.com").Data)
	assert.Empty(t, repo.GetAll("brainbaking.com").Data)
}

func TestHandleApprove(t *testing.T) {
	wm := mf.Mention{
		Source: "https://infos.by/markdown-v-nauke/",
		Target: "https://brainbaking.com/post/2021/02/writing-academic-papers-in-markdown/",
	}

	key, _ := repo.InModeration(wm, &mf.IndiewebData{
		Source: wm.Source,
		Target: wm.Target,
		Name:   "mytest",
	})
	assert.NotEmpty(t, repo.GetAllToModerate("brainbaking.com").Data)

	r := mux.NewRouter()
	// just using httptest.NewServer(r.HandleFunc("url", HandleApprove(...)) won't get the context vars into the mux
	r.HandleFunc("/admin/approve/{token}/{key}", HandleApprove(cnf, repo)).Methods("GET")
	ts := httptest.NewServer(r)

	t.Cleanup(func() {
		os.Remove("config.json")
		ts.Close()
		db.Purge()
	})

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/admin/approve/%s/%s", ts.URL, cnf.Token, key), nil)

	_, err = client.Do(req)
	assert.NoError(t, err)

	assert.Empty(t, repo.GetAllToModerate("brainbaking.com").Data)
	assert.NotEmpty(t, repo.GetAll("brainbaking.com").Data)
}
