package external

import (
	"brainbaking.com/go-jamming/app/mf"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTryImportBridgyUrl(t *testing.T) {
	wmio := &WebmentionIOImporter{}
	cases := []struct {
		label          string
		mention        string
		expectedSource string
	}{
		{
			"conventional source URL does nothing special",
			`{ "links": [ { "source": "https://brainbaking.com/lolz"  } ] }`,
			"https://brainbaking.com/lolz",
		},
		{
			"Source URL from brid.gy takes data URL as source instead",
			`{ "links": [ { "source": "https://brid.gy/like/twitter/iamchrisburnell/1298550501307486208/252048752", "data": { "url": "https://twitter.com/iamchrisburnell/status/1298550501307486208#favorited-by-252048752" } } ] }`,
			"https://twitter.com/iamchrisburnell/status/1298550501307486208#favorited-by-252048752",
		},
		{
			"Source URL from brid-gy.appspot.com takes URL as data source instead",
			`{ "links": [ { "source": "https://brid-gy.appspot.com/post/twitter/iamchrisburnell/1103728693648809984", "data": { "url": "https://twitter.com/adactioLinks/status/1103728693648809984" } } ] }`,
			"https://twitter.com/adactioLinks/status/1103728693648809984",
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			res, err := wmio.TryImport([]byte(tc.mention))
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedSource, res[0].Source)
		})
	}
}

func TestTryImportPublishedDates(t *testing.T) {
	wmio := &WebmentionIOImporter{}
	cases := []struct {
		label        string
		mention      string
		expectedDate string
	}{
		{
			"no dates reverts to first",
			`{ "links": [ {   } ] }`,
			time.Time{}.Format(mf.DateFormatWithTimeZone),
		},
		{
			"no published date reverts to verified date",
			`{ "links": [ { "verified_date": "2022-05-25T14:28:10+00:00"  } ] }`,
			"2022-05-25T14:28:10+00:00",
		},
		{
			"published date present takes preference over rest",
			`{ "links": [ { "data": { "published": "2020-01-25T14:28:10+00:00" }, "verified_date": "2022-05-25T14:28:10+00:00"  } ] }`,
			"2020-01-25T14:28:10+00:00",
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			res, err := wmio.TryImport([]byte(tc.mention))
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedDate, res[0].Published)
		})
	}
}

func TestTryImportErrorIfInvalidFormat(t *testing.T) {
	wmio := &WebmentionIOImporter{}
	mention := `haha`

	_, err := wmio.TryImport([]byte(mention))
	assert.Error(t, err)
}

func TestTryImportForLikeWithMissingAuthor(t *testing.T) {
	wmio := &WebmentionIOImporter{}
	mention := `{ "links": [
{
            "source": "https://jacky.wtf/2022/5/BRQo",
            "verified": true,
            "verified_date": "2022-05-25T14:28:10+00:00",
            "id": 1404286,
            "private": false,
            "data": {
                "url": "https://jacky.wtf/2022/5/BRQo",
                "name": null,
                "content": null,
                "published": "2022-05-25T14:26:12+00:00",
                "published_ts": 1653488772
            },
            "activity": {
                "type": "like",
                "sentence": "https://jacky.wtf/2022/5/BRQo liked a post https://chrisburnell.com/article/changing-with-the-times/",
                "sentence_html": "<a href=\"https://jacky.wtf/2022/5/BRQo\">someone</a> liked a post <a href=\"https://chrisburnell.com/article/changing-with-the-times/\">https://chrisburnell.com/article/changing-with-the-times/</a>"
            },
            "target": "https://chrisburnell.com/article/changing-with-the-times/"
        }
] }`

	res, err := wmio.TryImport([]byte(mention))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	result := res[0]

	assert.Equal(t, "https://chrisburnell.com/article/changing-with-the-times/", result.Target)
	assert.Equal(t, "https://jacky.wtf/2022/5/BRQo", result.Source)

	assert.Equal(t, mf.TypeLike, result.IndiewebType)
	assert.Equal(t, "https://jacky.wtf/2022/5/BRQo liked a post https://chrisburnell.com/article/changing-with-the-times/", result.Content)
	assert.Equal(t, "", result.Name)
	assert.Equal(t, "https://jacky.wtf/2022/5/BRQo", result.Url)
	assert.Equal(t, "2022-05-25T14:26:12+00:00", result.Published)

	assert.Equal(t, "", result.Author.Name)
	assert.Equal(t, "", result.Author.Picture)
}

func TestTryImportForReply(t *testing.T) {
	wmio := &WebmentionIOImporter{}
	mention := `{ "links": [
	{
		"source": "https://chrisburnell.com/note/1652693080/",
		"verified": true,
		"verified_date": "2022-05-16T19:36:52+00:00",
		"id": 1399408,
		"private": false,
		"data": {
		"author": {
			"name": "Chris Burnell",
				"url": "https://chrisburnell.com/",
				"photo": "https://webmention.io/avatar/chrisburnell.com/ace41559b8d4e8d8189b285d88b1ea2dc6c53056fc512be7d199c0c8cadc53fe.jpg"
		},
		"url": "https://chrisburnell.com/note/1652693080/",
			"name": null,
			"content": "<p>first!!1!</p>",
			"published": "2022-05-16T09:24:40+01:00",
			"published_ts": 1652689480
	},
		"activity": {
		"type": "reply",
			"sentence": "Chris Burnell commented 'first!!1!' on a post https://chrisburnell.com/guestbook/",
			"sentence_html": "<a href=\"https://chrisburnell.com/\">Chris Burnell</a> commented 'first!!1!' on a post <a href=\"https://chrisburnell.com/guestbook/\">https://chrisburnell.com/guestbook/</a>"
	},
		"rels": {
		"canonical": "https://chrisburnell.com/note/1652693080/"
	},
		"target": "https://chrisburnell.com/guestbook/"
	}
] }`

	res, err := wmio.TryImport([]byte(mention))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	result := res[0]

	assert.Equal(t, "https://chrisburnell.com/guestbook/", result.Target)
	assert.Equal(t, "https://chrisburnell.com/note/1652693080/", result.Source)

	assert.Equal(t, mf.TypeReply, result.IndiewebType)
	assert.Equal(t, "Chris Burnell commented 'first!!1!' on a post https://chrisburnell.com/guestbook/", result.Content)
	assert.Equal(t, "<p>first!!1!</p>", result.Name)
	assert.Equal(t, "https://chrisburnell.com/note/1652693080/", result.Url)
	assert.Equal(t, "2022-05-16T09:24:40+01:00", result.Published)

	assert.Equal(t, "Chris Burnell", result.Author.Name)
	assert.Equal(t, "https://webmention.io/avatar/chrisburnell.com/ace41559b8d4e8d8189b285d88b1ea2dc6c53056fc512be7d199c0c8cadc53fe.jpg", result.Author.Picture)
}
