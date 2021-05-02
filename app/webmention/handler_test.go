package webmention

import (
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
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
	cnf  = common.Configure()
	repo db.MentionRepo
)

func init() {
	cnf.ConString = ":memory:"
	repo = db.NewMentionRepo(cnf)
}

func TestHandlePostWithInvalidUrlsShouldReturnBadRequest(t *testing.T) {
	ts := httptest.NewServer(HandlePost(cnf, repo))
	defer ts.Close()
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
