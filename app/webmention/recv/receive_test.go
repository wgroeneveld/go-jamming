
package recv

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"brainbaking.com/go-jamming/app/mf"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/mocks"
)

var conf = &common.Config{
	AllowedWebmentionSources: []string {
		"jefklakscodex.com",
	},
	DataPath: "testdata",
}


func TestConvertWebmentionToPath(t *testing.T) {
	wm := mf.Mention{
		Source: "https://brainbaking.com",
		Target: "https://jefklakscodex.com/articles",
	}

	result := wm.AsPath(conf)
	if result != "testdata/jefklakscodex.com/99be66594fdfcf482545fead8e7e4948.json" {
		t.Fatalf("md5 hash check failed, got " + result)
	}
}

func writeSomethingTo(filename string) {
	file, _ := os.Create(filename)
	file.WriteString("lolz")
	defer file.Close()	
}

func TestReceive(t *testing.T) {
	cases := []struct {
		label string
		wm    mf.Mention
		json  string
	} {
		{
			label: "receive a Webmention bookmark via twitter",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-bridgy-twitter-source.html",
				Target: "https://brainbaking.com/post/2021/03/the-indieweb-mixed-bag",
			},
			json: `{"author":{"name":"Jamie Tanna","picture":"https://www.jvt.me/img/profile.png"},"name":"","content":"Recommended read:\nThe IndieWeb Mixed Bag - Thoughts about the (d)evolution of blog interactions\nhttps://brainbaking.com/post/2021/03/the-indieweb-mixed-bag/","published":"2021-03-15T12:42:00+0000","url":"https://brainbaking.com/mf2/2021/03/1bkre/","type":"bookmark","source":"https://brainbaking.com/valid-bridgy-twitter-source.html","target":"https://brainbaking.com/post/2021/03/the-indieweb-mixed-bag"}`,
		},
		{
			label: "receive a brid.gy Webmention like",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-bridgy-like.html",
				// wrapped in a a class="u-like-of" tag
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			// no dates in bridgy-to-mastodon likes...
			json: `{"author":{"name":"Stampeding Longhorn","picture":"https://cdn.social.linux.pizza/v1/AUTH_91eb37814936490c95da7b85993cc2ff/sociallinuxpizza/accounts/avatars/000/185/996/original/9e36da0c093cfc9b.png"},"name":"","content":"","published":"2020-01-01T12:30:00","url":"https://chat.brainbaking.com/notice/A4nx1rFwKUJYSe4TqK#favorited-by-A4nwg4LYyh4WgrJOXg","type":"like","source":"https://brainbaking.com/valid-bridgy-like.html","target":"https://brainbaking.com/valid-indieweb-target.html"}`,
		},
		{
			label: "receive a brid.gy Webmention that has a url and photo without value",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-bridgy-source.html",
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			json: `{"author":{"name":"Stampeding Longhorn", "picture":"https://cdn.social.linux.pizza/v1/AUTH_91eb37814936490c95da7b85993cc2ff/sociallinuxpizza/accounts/avatars/000/185/996/original/9e36da0c093cfc9b.png"}, "content":"@wouter The cat pictures are awesome. for jest tests!", "name":"@wouter The cat pictures are awesome. for jest tests!", "published":"2021-03-02T16:17:18.000Z", "source":"https://brainbaking.com/valid-bridgy-source.html", "target":"https://brainbaking.com/valid-indieweb-target.html", "type":"mention", "url":"https://social.linux.pizza/@StampedingLonghorn/105821099684887793"}`,
		},
		{
			label: "receive saves a JSON file of indieweb-metadata if all is valid",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-indieweb-source.html",
				Target: "https://jefklakscodex.com/articles",
			},
			json: `{"author":{"name":"Wouter Groeneveld","picture":"https://brainbaking.com//img/avatar.jpg"},"name":"I just learned about https://www.inklestudios.com/...","content":"This is cool, I just found out about valid indieweb target - so cool","published":"2021-03-06T12:41:00","url":"https://brainbaking.com/notes/2021/03/06h12m41s48/","type":"mention","source":"https://brainbaking.com/valid-indieweb-source.html","target":"https://jefklakscodex.com/articles"}`,
		},
		{
			label: "receive saves a JSON file of indieweb-metadata with summary as content if present",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-indieweb-source-with-summary.html",
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			json: `{"author":{"name":"Wouter Groeneveld", "picture":"https://brainbaking.com//img/avatar.jpg"}, "content":"This is cool, this is a summary!", "name":"I just learned about https://www.inklestudios.com/...", "published":"2021-03-06T12:41:00", "source":"https://brainbaking.com/valid-indieweb-source-with-summary.html", "target":"https://brainbaking.com/valid-indieweb-target.html", "type":"mention", "url":"https://brainbaking.com/notes/2021/03/06h12m41s48/"}`,
		},
		{
			label: "receive saves a JSON file of non-indieweb-data such as title if all is valid",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-nonindieweb-source.html",
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			json: `{"author":{"name":"https://brainbaking.com/valid-nonindieweb-source.html", "picture":""}, "content":"Diablo 2 Twenty Years Later: A Retrospective | Jefklaks Codex", "name":"Diablo 2 Twenty Years Later: A Retrospective | Jefklaks Codex", "published":"2020-01-01T12:30:00", "source":"https://brainbaking.com/valid-nonindieweb-source.html", "target":"https://brainbaking.com/valid-indieweb-target.html", "type":"mention", "url":"https://brainbaking.com/valid-nonindieweb-source.html"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			os.MkdirAll("testdata/brainbaking.com", os.ModePerm)
			os.MkdirAll("testdata/jefklakscodex.com", os.ModePerm)
			defer os.RemoveAll("testdata")
			common.Now = func() time.Time {
				return time.Date(2020, time.January, 1, 12, 30, 0, 0, time.UTC)
			}

			receiver := &Receiver{
				Conf: conf,
				RestClient: &mocks.RestClientMock{
					GetBodyFunc: mocks.RelPathGetBodyFunc(t, "../../../mocks/"),
				},
			}

			receiver.Receive(tc.wm)

			actualJson, _ := ioutil.ReadFile(tc.wm.AsPath(conf))
			assert.JSONEq(t, tc.json, string(actualJson))
		})
	}
}

func TestReceiveTargetDoesNotExistAnymoreDeletesPossiblyOlderWebmention(t *testing.T) {
	os.MkdirAll("testdata/jefklakscodex.com", os.ModePerm)
	defer os.RemoveAll("testdata")

	wm := mf.Mention{
		Source: "https://brainbaking.com",
		Target: "https://jefklakscodex.com/articles",
	}
	filename := wm.AsPath(conf)
	writeSomethingTo(filename)

	client := &mocks.RestClientMock{
		GetBodyFunc: func(url string) (string, error) {
			return "", errors.New("whoops")
		},
	}	
	receiver := &Receiver{
		Conf:       conf,
		RestClient: client,
	}

	receiver.Receive(wm)
	assert.NoFileExists(t, filename)
}

func TestReceiveTargetThatDoesNotPointToTheSourceDoesNothing(t *testing.T) {
	wm := mf.Mention{
		Source: "https://brainbaking.com/valid-indieweb-source.html",
		Target: "https://brainbaking.com/valid-indieweb-source.html",
	}
	filename := wm.AsPath(conf)
	writeSomethingTo(filename)

	receiver := &Receiver{
		Conf: conf,
		RestClient: &mocks.RestClientMock{
			GetBodyFunc: mocks.RelPathGetBodyFunc(t, "../../../mocks/"),
		},
	}

	receiver.Receive(wm)
	assert.NoFileExists(t, filename)
}

func TestProcessSourceBodyAbortsIfNoMentionOfTargetFoundInSourceHtml(t *testing.T) {
	os.MkdirAll("testdata/jefklakscodex.com", os.ModePerm)
	defer os.RemoveAll("testdata")

	wm := mf.Mention{
		Source: "https://brainbaking.com",
		Target: "https://jefklakscodex.com/articles",
	}
	receiver := &Receiver{
		Conf: conf,
	}

	receiver.processSourceBody("<html>my nice body</html>", wm)
	assert.NoFileExists(t, wm.AsPath(conf))
}

