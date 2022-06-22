package external

import (
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/mocks"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"testing"
)

var (
	cnf = &common.Config{
		BaseURL:                  "http://localhost:1337/",
		Port:                     1337,
		Token:                    "miauwkes",
		AllowedWebmentionSources: []string{"chrisburnell.com"},
		Blacklist:                []string{},
		Whitelist:                []string{"chrisburnell.com"},
	}
)

func TestImport(t *testing.T) {
	repo := db.NewMentionRepo(cnf)
	bootstrapper := ImportBootstrapper{
		Conf: cnf,
		Repo: repo,
		RestClient: &mocks.RestClientMock{
			// this will make sure each picture GET fails
			// otherwise this test is REALLY slow. It will fallback to anonymous pictures
			GetBodyFunc: mocks.RelPathGetBodyFunc("../../../mocks/"),
		},
	}

	t.Cleanup(func() {
		os.Remove("config.json")
		db.Purge()
	})

	bootstrapper.Import("../../mocks/external/wmio.json")

	entries := repo.GetAll("chrisburnell.com")
	assert.Equal(t, 25, len(entries.Data))
	sort.Slice(entries.Data, func(i, j int) bool {
		return entries.Data[i].PublishedDate().After(entries.Data[j].PublishedDate())
	})

	assert.Equal(t, "https://chrisburnell.com/note/1655219889/", entries.Data[0].Source)
	assert.Equal(t, "/pictures/anonymous", entries.Data[0].Author.Picture)
	assert.Equal(t, "", entries.Data[10].Name)
	assert.Equal(t, "https://jacky.wtf/2022/5/BRQo liked a post https://chrisburnell.com/article/changing-with-the-times/", entries.Data[20].Content)

	assert.Contains(t, cnf.Whitelist, "jacky.wtf")
	assert.Contains(t, cnf.Whitelist, "martymcgui.re")
}
