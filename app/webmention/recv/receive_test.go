package recv

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/notifier"
	"brainbaking.com/go-jamming/db"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
	"willnorris.com/go/microformats"

	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/mocks"
)

var conf = &common.Config{
	AllowedWebmentionSources: []string{
		"jefklakscodex.com",
		"brainbaking.com",
	},
	Blacklist: []string{
		"blacklisted.com",
	},
	Whitelist: []string{
		"brainbaking.com",
		"jefklakscodex.com",
	},
}

func TestSaveAuthorPictureLocally(t *testing.T) {
	cases := []struct {
		label              string
		pictureUrl         string
		expectedPictureUrl string
		expectedError      error
	}{
		{
			"Absolute URL gets 'downloaded' and replaced by relative",
			"https://brainbaking.com/picture.jpg",
			"/pictures/brainbaking.com",
			nil,
		},
		{
			"Refuses to download if it's from a silo domain and possibly involves GDPR privacy issues",
			"https://brid.gy/picture.jpg",
			"https://brid.gy/picture.jpg",
			errWontDownloadBecauseOfPrivacy,
		},
		{
			"Absolute URL does not get replaced but error if no valid image",
			"https://brainbaking.com/index.xml",
			"https://brainbaking.com/index.xml",
			errPicNoRealImage,
		},
		{
			"Absolute URL does not get replaced but error if download fails",
			"https://brainbaking.com/thedogatemypic-nowitsmissing-shiii.png",
			"https://brainbaking.com/thedogatemypic-nowitsmissing-shiii.png",
			errPicUnableToDownload,
		},
	}

	t.Cleanup(db.Purge)
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			repo := db.NewMentionRepo(conf)
			recv := &Receiver{
				Conf: conf,
				Repo: repo,
				RestClient: &mocks.RestClientMock{
					GetBodyFunc: mocks.RelPathGetBodyFunc("../../../mocks/"),
				},
			}

			indieweb := &mf.IndiewebData{
				Source: tc.pictureUrl,
				Author: mf.IndiewebAuthor{
					Picture: tc.pictureUrl,
				},
			}
			err := recv.saveAuthorPictureLocally(indieweb)

			assert.Equal(t, tc.expectedPictureUrl, indieweb.Author.Picture)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func BenchmarkReceiveWithoutRestCalls(b *testing.B) {
	origLog := zerolog.GlobalLevel()
	zerolog.SetGlobalLevel(zerolog.Disabled)
	defer zerolog.SetGlobalLevel(origLog)

	wm := mf.Mention{
		Source: "https://brainbaking.com/valid-indieweb-source.html",
		Target: "https://jefklakscodex.com/articles",
	}
	data, err := ioutil.ReadFile("../../../mocks/valid-indieweb-source.html")
	assert.NoError(b, err)
	html := string(data)

	repo := db.NewMentionRepo(conf)
	recv := &Receiver{
		Conf: conf,
		Repo: repo,
		RestClient: &mocks.RestClientMock{
			GetBodyFunc: func(s string) (http.Header, string, error) {
				return http.Header{}, html, nil
			},
		},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		recv.Receive(wm)
	}
}

func TestReceive(t *testing.T) {
	cases := []struct {
		label string
		wm    mf.Mention
		json  string
	}{
		{
			label: "bugfix interface conversion panic unusual author part in mf",
			wm: mf.Mention{
				Source: "https://brainbaking.com/bugfix-interface-conversion-panic.html",
				Target: "https://brainbaking.com/",
			},
			json: `{"author":{"name":"Ton Zijlstra","picture":"/pictures/brainbaking.com"},"name":"","content":"De allereerste Nederlandstalige meet-up van Obsidian.md gebruikers was interessant en leuk! We waren met z’n vieren, Sebastiaan, Wouter, Frank en ik, en spraken bijna 2 uur met elkaar. Leuk om te vergelijken waarom en hoe we notities maken in Obsid...", "name":"Nabeschouwing: de eerste Nederlandstalige Obsidian meet-up", "published":"2021-04-25T11:24:48+02:00", "source":"https://brainbaking.com/bugfix-interface-conversion-panic.html", "target":"https://brainbaking.com/", "type":"mention", "url":"https://www.zylstra.org/blog/2021/04/nabeschouwing-de-eerste-nederlandstalige-obsidian-meet-up/"}`,
		},
		{
			label: "receive a Webmention bookmark via twitter",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-bridgy-twitter-source.html",
				Target: "https://brainbaking.com/post/2021/03/the-indieweb-mixed-bag",
			},
			json: `{"author":{"name":"Jamie Tanna","picture":"/pictures/brainbaking.com"},"name":"","content":"Recommended read:\nThe IndieWeb Mixed Bag - Thoughts about the (d)evolution of blog interactions\nhttps://brainbaking.com/post/2021/03/the-indieweb-mixed-bag/","published":"2021-03-15T12:42:00+00:00","url":"https://brainbaking.com/mf2/2021/03/1bkre/","type":"bookmark","source":"https://brainbaking.com/valid-bridgy-twitter-source.html","target":"https://brainbaking.com/post/2021/03/the-indieweb-mixed-bag"}`,
		},
		{
			label: "receive a brid.gy (Mastodon) Webmention like",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-bridgy-like.html",
				// wrapped in a a class="u-like-of" tag
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			// no dates in bridgy-to-mastodon likes...
			json: `{"author":{"name":"Stampeding Longhorn","picture":"/pictures/brainbaking.com"},"name":"","content":"","published":"2020-01-01T12:30:00+00:00","url":"https://chat.brainbaking.com/notice/A4nx1rFwKUJYSe4TqK#favorited-by-A4nwg4LYyh4WgrJOXg","type":"like","source":"https://brainbaking.com/valid-bridgy-like.html","target":"https://brainbaking.com/valid-indieweb-target.html"}`,
		},
		{
			label: "receive a brid.gy (Twitter) Webmention repost",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-bridgy-twitter-repost.html",
				// wrapped in a a class="u-like-of" tag
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			// no dates in bridgy-to-mastodon likes...
			json: `{"author":{"name":"cartocalypse.tif","picture":"/pictures/brainbaking.com"},"name":"My quest in creating a Google Maps clone\n\n            chringel.dev/2022/06/creati…","content":"My quest in creating a Google Maps clone\n\n            chringel.dev/2022/06/creati…","published":"2022-06-21T06:23:53+00:00","url":"https://twitter.com/cartocalypse/status/1539131976879308800","type":"repost","source":"https://brainbaking.com/valid-bridgy-twitter-repost.html","target":"https://brainbaking.com/valid-indieweb-target.html"}`,
		},
		{
			label: "receive a brid.gy (Mastodon) Webmention that has a url and photo without value",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-bridgy-source.html",
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			json: `{"author":{"name":"Stampeding Longhorn", "picture":"/pictures/brainbaking.com"}, "content":"@wouter The cat pictures are awesome. for jest tests!", "name":"@wouter The cat pictures are awesome. for jest tests!", "published":"2021-03-02T16:17:18+00:00", "source":"https://brainbaking.com/valid-bridgy-source.html", "target":"https://brainbaking.com/valid-indieweb-target.html", "type":"mention", "url":"https://social.linux.pizza/@StampedingLonghorn/105821099684887793"}`,
		},
		{
			label: "receive saves a JSON file of indieweb-metadata if all is valid",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-indieweb-source.html",
				Target: "https://jefklakscodex.com/articles",
			},
			json: `{"author":{"name":"Wouter Groeneveld","picture":"/pictures/brainbaking.com"},"name":"I just learned about https://www.inklestudios.com/...","content":"This is cool, I just found out about valid indieweb target - so cool","published":"2021-03-06T12:41:00+00:00","url":"https://brainbaking.com/notes/2021/03/06h12m41s48/","type":"mention","source":"https://brainbaking.com/valid-indieweb-source.html","target":"https://jefklakscodex.com/articles"}`,
		},
		{
			label: "receive saves a JSON file of indieweb-metadata with summary as content if present",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-indieweb-source-with-summary.html",
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			json: `{"author":{"name":"Wouter Groeneveld", "picture":"/pictures/brainbaking.com"}, "content":"This is cool, this is a summary!", "name":"I just learned about https://www.inklestudios.com/...", "published":"2021-03-06T12:41:00+00:00", "source":"https://brainbaking.com/valid-indieweb-source-with-summary.html", "target":"https://brainbaking.com/valid-indieweb-target.html", "type":"mention", "url":"https://brainbaking.com/notes/2021/03/06h12m41s48/"}`,
		},
		{
			label: "receive saves a JSON file of non-indieweb-data such as title if all is valid",
			wm: mf.Mention{
				Source: "https://brainbaking.com/valid-nonindieweb-source.html",
				Target: "https://brainbaking.com/valid-indieweb-target.html",
			},
			json: `{"author":{"name":"https://brainbaking.com/valid-nonindieweb-source.html", "picture":"/pictures/anonymous"}, "content":"Diablo 2 Twenty Years Later: A Retrospective | Jefklaks Codex", "name":"Diablo 2 Twenty Years Later: A Retrospective | Jefklaks Codex", "published":"2020-01-01T12:30:00+00:00", "source":"https://brainbaking.com/valid-nonindieweb-source.html", "target":"https://brainbaking.com/valid-indieweb-target.html", "type":"mention", "url":"https://brainbaking.com/valid-nonindieweb-source.html"}`,
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
					GetBodyFunc: mocks.RelPathGetBodyFunc("../../../mocks/"),
				},
				Notifier: &notifier.StringNotifier{},
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
	t.Cleanup(db.Purge)

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

func TestReceiveFromNotInWhitelistSavesInModerationAndNotifies(t *testing.T) {
	wm := mf.Mention{
		Source: "https://brainbaking.com/valid-indieweb-source.html",
		Target: "https://brainbaking.com/valid-indieweb-target.html",
	}
	cnf := &common.Config{
		AllowedWebmentionSources: []string{
			"brainbaking.com",
		},
		BaseURL:   "https://jam.brainbaking.com/",
		Token:     "mytoken",
		Blacklist: []string{},
		Whitelist: []string{},
	}
	repo := db.NewMentionRepo(cnf)
	t.Cleanup(db.Purge)
	notifierMock := &notifier.StringNotifier{
		Conf:   cnf,
		Output: "",
	}
	receiver := &Receiver{
		Conf: cnf,
		Repo: repo,
		RestClient: &mocks.RestClientMock{
			GetBodyFunc: mocks.RelPathGetBodyFunc("../../../mocks/"),
		},
		Notifier: notifierMock,
	}

	receiver.Receive(wm)
	assert.Empty(t, repo.GetAll("brainbaking.com").Data)
	assert.Equal(t, 1, len(repo.GetAllToModerate("brainbaking.com").Data))
	assert.Contains(t, notifierMock.Output, "in moderation!")
}

func TestReceiveFromBlacklistedDomainDoesNothing(t *testing.T) {
	wm := mf.Mention{
		Source: "https://blacklisted.com/whoops",
		Target: "https://brainbaking.com/valid-indieweb-source.html",
	}

	repo := db.NewMentionRepo(conf)
	t.Cleanup(db.Purge)
	receiver := &Receiver{
		Conf: conf,
		Repo: repo,
	}

	receiver.Receive(wm)
	assert.Empty(t, repo.GetAll("brainbaking.com").Data)
	assert.Empty(t, repo.GetAllToModerate("brainbaking.com").Data)
}

func TestReceiveTargetThatDoesNotPointToTheSourceDoesNothing(t *testing.T) {
	wm := mf.Mention{
		Source: "https://brainbaking.com/valid-indieweb-source.html",
		Target: "https://brainbaking.com/valid-indieweb-source.html",
	}

	repo := db.NewMentionRepo(conf)
	t.Cleanup(db.Purge)
	receiver := &Receiver{
		Conf: conf,
		Repo: repo,
		RestClient: &mocks.RestClientMock{
			GetBodyFunc: mocks.RelPathGetBodyFunc("../../../mocks/"),
		},
	}

	receiver.Receive(wm)
	assert.Empty(t, repo.GetAll("brainbaking.com").Data)
	assert.Empty(t, repo.GetAllToModerate("brainbaking.com").Data)
}

func TestConvertBodyToIndiewebDataWithComplicatedData(t *testing.T) {
	cases := []struct {
		source             string
		target             string
		mockFile           string
		expectedAuthorName string
	}{
		{
			"https://pauho.net/2022/05/04/662325/",
			"https://brainbaking.com/post/2022/04/cool-things-people-do-with-their-blogs/",
			"../../../mocks/indieweb-complicated-data.html",
			"Paul Houlihan",
		},
		{
			"https://ruk.ca/favourite/drought",
			"https://brainbaking.com/post/2022/05/drought/",
			"../../../mocks/indieweb-hcard-next-to-name.html",
			"Peter Rukavina",
		},
	}

	cnf := &common.Config{
		AllowedWebmentionSources: []string{
			"brainbaking.com",
		},
	}
	recv := &Receiver{
		Conf: cnf,
	}
	for _, tc := range cases {
		t.Run("indieweb author "+tc.expectedAuthorName, func(t *testing.T) {
			wm := mf.Mention{
				Source: tc.source,
				Target: tc.target,
			}

			src, err := ioutil.ReadFile(tc.mockFile)
			assert.NoError(t, err)

			body := string(src)
			data := microformats.Parse(strings.NewReader(body), wm.SourceUrl())
			indieweb := recv.convertBodyToIndiewebData(body, wm, data)

			assert.Equal(t, tc.expectedAuthorName, indieweb.Author.Name)
		})
	}
}

func TestProcessSourceBodyAnonymizesBothAuthorPictureAndNameIfComingFromSilo(t *testing.T) {
	wm := mf.Mention{
		Source: "https://brid.gy/post/twitter/ChrisAldrich/1387130900962443264",
		Target: "https://brainbaking.com/",
	}
	cnf := &common.Config{
		AllowedWebmentionSources: []string{
			"brainbaking.com",
		},
		Whitelist: []string{
			"brid.gy",
		},
	}

	repo := db.NewMentionRepo(cnf)
	t.Cleanup(db.Purge)
	recv := &Receiver{
		Conf: cnf,
		Repo: repo,
	}

	src, err := ioutil.ReadFile("../../../mocks/valid-bridgy-source.html")
	assert.NoError(t, err)

	recv.processSourceBody(string(src), wm)
	savedMention := repo.Get(wm)

	assert.Equal(t, "Anonymous", savedMention.Author.Name)
	assert.Equal(t, "/pictures/anonymous", savedMention.Author.Picture)
}

func TestProcessSourceBodyAbortsIfNoMentionOfTargetFoundInSourceHtml(t *testing.T) {
	wm := mf.Mention{
		Source: "https://brainbaking.com",
		Target: "https://jefklakscodex.com/articles",
	}
	repo := db.NewMentionRepo(conf)
	t.Cleanup(db.Purge)
	recv := &Receiver{
		Conf: conf,
		Repo: repo,
	}

	recv.processSourceBody("<html>my nice body</html>", wm)
	assert.Empty(t, repo.Get(wm))
}

func TestProcessAuthorPictureAnonymizesIfEmpty(t *testing.T) {
	recv := &Receiver{}
	indieweb := &mf.IndiewebData{
		Author: mf.IndiewebAuthor{
			Picture: "",
		},
	}
	recv.ProcessAuthorPicture(indieweb)

	assert.Equal(t, "/pictures/anonymous", indieweb.Author.Picture)
}
