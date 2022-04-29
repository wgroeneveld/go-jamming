package webmention

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

func postWm(source string, target string) url.Values {
	values := url.Values{}
	values.Set("source", source)
	values.Set("target", target)
	return values
}

var (
	cnf = &common.Config{
		BaseURL:                  "http://localhost:1337/",
		Port:                     1337,
		Token:                    "miauwkes",
		AllowedWebmentionSources: []string{"brainbaking.com"},
		Blacklist:                []string{"youtube.com"},
		Whitelist:                []string{"brainbaking.com"},
	}
	repo db.MentionRepo
)

func init() {
	repo = db.NewMentionRepo(cnf)
}

func TestHandleDelete(t *testing.T) {
	wm := mf.Mention{
		Source: "https://infos.by/markdown-v-nauke/",
		Target: "https://brainbaking.com/post/2021/02/writing-academic-papers-in-markdown/",
	}

	_, err := repo.Save(wm, &mf.IndiewebData{
		Source: wm.Source,
		Target: wm.Target,
		Name:   "mytest",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, repo.GetAll("brainbaking.com").Data)

	ts := httptest.NewServer(HandleDelete(repo))
	defer ts.Close()
	t.Cleanup(db.Purge)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s?source=%s&target=%s", ts.URL, wm.Source, wm.Target), nil)
	assert.NoError(t, err)
	_, err = client.Do(req)
	assert.NoError(t, err)

	assert.Empty(t, repo.GetAll("brainbaking.com").Data)
}

func TestHandlePostWithInvalidUrlsShouldReturnBadRequest(t *testing.T) {
	ts := httptest.NewServer(HandlePost(cnf, repo))
	defer ts.Close()
	t.Cleanup(db.Purge)

	res, err := http.PostForm(ts.URL, postWm("https://haha.be/woof/said/the/dog.txt", "https://pussies.nl/mycatjustthrewup/gottacleanup.html"))
	assert.NoError(t, err)

	content, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	assert.NoError(t, err)
	assert.Contains(t, string(content), "Bad Request")
}

// Explicitly tests using actual live data, so this could fail if URLs are unreachable.
func TestHandlePostWithTestServer_Parallel(t *testing.T) {
	ts := httptest.NewServer(HandlePost(cnf, repo))
	defer ts.Close()
	t.Cleanup(db.Purge)

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := http.PostForm(ts.URL, postWm("https://jefklakscodex.com/articles/retrospectives/raven-shield-17-years-later/", "https://brainbaking.com/post/2020/10/building-a-core2duo-winxp-retro-pc/"))
			assert.NoError(t, err)

			content, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			assert.NoError(t, err)
			assert.Contains(t, string(content), "Thanks, bro. Will process this soon, pinky swear")
		}()
	}
	wg.Wait()
}
