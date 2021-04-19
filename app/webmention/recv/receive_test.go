package recv

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/db"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"

	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/mocks"
)

var conf = &common.Config{
	AllowedWebmentionSources: []string{
		"jefklakscodex.com",
		"brainbaking.com",
	},
	ConString: ":memory:",
}

func TestSaveAuthorPictureLocally(t *testing.T) {
	cases := []struct {
		label              string
		pictureUrl         string
		expectedPictureUrl string
	}{
		{
			"Absolute URL gets 'downloaded' and replaced by relative",
			"https://brainbaking.com/picture.jpg",
			"/pictures/brainbaking.com",
		},
		{
			"Absolute URL gets replaced by anonymous if download fails",
			"https://brainbaking.com/thedogatemypic-nowitsmissing-shiii.png",
			"/pictures/anonymous",
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			repo := db.NewMentionRepo(conf)
			recv := &Receiver{
				Conf: conf,
				Repo: repo,
				RestClient: &mocks.RestClientMock{
					GetBodyFunc: mocks.RelPathGetBodyFunc(t, "../../../mocks/"),
				},
			}

			indieweb := &mf.IndiewebData{
				Source: "https://brainbaking.com",
				Author: mf.IndiewebAuthor{
					Picture: tc.pictureUrl,
				},
			}
			recv.saveAuthorPictureLocally(indieweb)

			assert.Equal(t, tc.expectedPictureUrl, indieweb.Author.Picture)
		})
	}
}

func TestReceive(t *testing.T) {
	cases := []struct {
		label string
		wm    mf.Mention
		json  string
	}{
		{
			label: "receive a Webmention bookmark via twitter",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-bridgy-twitter-source.html",
				Target: "https://brainbaking.com/post/2021/03/the-indieweb-mixed-bag",
			},
			json: `{"author":{"name":"Jamie Tanna","picture":"/pictures/brainbaking.com"},"name":"","content":"Recommended read:\nThe IndieWeb Mixed Bag - Thoughts about the (d)evolution of blog interactions\nhttps://brainbaking.com/post/2021/03/the-indieweb-mixed-bag/","published":"2021-03-15T12:42:00+0000","url":"https://brainbaking.com/mf2/2021/03/1bkre/","type":"bookmark","source":"https://brainbaking.com/valid-bridgy-twitter-source.html","target":"https://brainbaking.com/post/2021/03/the-indieweb-mixed-bag"}`,
		},
		{
			label: "receive a brid.gy Webmention like",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-bridgy-like.html",
				// wrapped in a a class="u-like-of" tag
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			// no dates in bridgy-to-mastodon likes...
			json: `{"author":{"name":"Stampeding Longhorn","picture":"/pictures/brainbaking.com"},"name":"","content":"","published":"2020-01-01T12:30:00","url":"https://chat.brainbaking.com/notice/A4nx1rFwKUJYSe4TqK#favorited-by-A4nwg4LYyh4WgrJOXg","type":"like","source":"https://brainbaking.com/valid-bridgy-like.html","target":"https://brainbaking.com/valid-indieweb-target.html"}`,
		},
		{
			label: "receive a brid.gy Webmention that has a url and photo without value",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-bridgy-source.html",
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			json: `{"author":{"name":"Stampeding Longhorn", "picture":"/pictures/brainbaking.com"}, "content":"@wouter The cat pictures are awesome. for jest tests!", "name":"@wouter The cat pictures are awesome. for jest tests!", "published":"2021-03-02T16:17:18.000Z", "source":"https://brainbaking.com/valid-bridgy-source.html", "target":"https://brainbaking.com/valid-indieweb-target.html", "type":"mention", "url":"https://social.linux.pizza/@StampedingLonghorn/105821099684887793"}`,
		},
		{
			label: "receive saves a JSON file of indieweb-metadata if all is valid",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-indieweb-source.html",
				Target: "https://jefklakscodex.com/articles",
			},
			json: `{"author":{"name":"Wouter Groeneveld","picture":"/pictures/brainbaking.com"},"name":"I just learned about https://www.inklestudios.com/...","content":"This is cool, I just found out about valid indieweb target - so cool","published":"2021-03-06T12:41:00","url":"https://brainbaking.com/notes/2021/03/06h12m41s48/","type":"mention","source":"https://brainbaking.com/valid-indieweb-source.html","target":"https://jefklakscodex.com/articles"}`,
		},
		{
			label: "receive saves a JSON file of indieweb-metadata with summary as content if present",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-indieweb-source-with-summary.html",
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			json: `{"author":{"name":"Wouter Groeneveld", "picture":"/pictures/brainbaking.com"}, "content":"This is cool, this is a summary!", "name":"I just learned about https://www.inklestudios.com/...", "published":"2021-03-06T12:41:00", "source":"https://brainbaking.com/valid-indieweb-source-with-summary.html", "target":"https://brainbaking.com/valid-indieweb-target.html", "type":"mention", "url":"https://brainbaking.com/notes/2021/03/06h12m41s48/"}`,
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
			common.Now = func() time.Time {
				return time.Date(2020, time.January, 1, 12, 30, 0, 0, time.UTC)
			}

			repo := db.NewMentionRepo(conf)
			receiver := &Receiver{
				Conf: conf,
				Repo: repo,
				RestClient: &mocks.RestClientMock{
					GetBodyFunc: mocks.RelPathGetBodyFunc(t, "../../../mocks/"),
				},
			}

			receiver.Receive(tc.wm)

			actual := repo.Get(tc.wm)
			actualJson, _ := json.Marshal(actual)
			assert.JSONEq(t, tc.json, string(actualJson))
		})
	}
}

func TestReceiveTargetDoesNotExistAnymoreDeletesPossiblyOlderWebmention(t *testing.T) {
	repo := db.NewMentionRepo(conf)

	wm := mf.Mention{
		Source: "https://brainbaking.com",
		Target: "https://jefklakscodex.com/articles",
	}
	repo.Save(wm, &mf.IndiewebData{
		Name: "something something",
	})

	client := &mocks.RestClientMock{
		GetBodyFunc: func(url string) (http.Header, string, error) {
			return nil, "", errors.New("whoops")
		},
	}
	receiver := &Receiver{
		Conf:       conf,
		RestClient: client,
		Repo:       repo,
	}

	receiver.Receive(wm)
	indb := repo.Get(wm)
	assert.Empty(t, indb)
}

func TestReceiveTargetThatDoesNotPointToTheSourceDoesNothing(t *testing.T) {
	wm := mf.Mention{
		Source: "https://brainbaking.com/valid-indieweb-source.html",
		Target: "https://brainbaking.com/valid-indieweb-source.html",
	}

	repo := db.NewMentionRepo(conf)
	receiver := &Receiver{
		Conf: conf,
		Repo: repo,
		RestClient: &mocks.RestClientMock{
			GetBodyFunc: mocks.RelPathGetBodyFunc(t, "../../../mocks/"),
		},
	}

	receiver.Receive(wm)
	assert.Empty(t, repo.GetAll("brainbaking.com").Data)
}

func TestProcessSourceBodyAbortsIfNoMentionOfTargetFoundInSourceHtml(t *testing.T) {
	wm := mf.Mention{
		Source: "https://brainbaking.com",
		Target: "https://jefklakscodex.com/articles",
	}
	repo := db.NewMentionRepo(conf)
	receiver := &Receiver{
		Conf: conf,
		Repo: repo,
	}

	receiver.processSourceBody("<html>my nice body</html>", wm)
	assert.Empty(t, repo.Get(wm))
}
