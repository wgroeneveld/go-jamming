package external

import (
	"brainbaking.com/go-jamming/app/mf"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
