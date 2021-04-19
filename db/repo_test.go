package db

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/buntdb"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var (
	conf = &common.Config{
		ConString: ":memory:",
		AllowedWebmentionSources: []string{
			"pussycat.com",
		},
	}
)

func TestSaveAndGetPicture(t *testing.T) {
	data, err := ioutil.ReadFile("../mocks/picture.jpg")
	assert.NoError(t, err)

	db := NewMentionRepo(conf)
	key, dberr := db.SavePicture(string(data), "bloeberig.be")
	assert.NoError(t, dberr)
	assert.Equal(t, "bloeberig.be:picture", key)

	picDataAfterSave := db.GetPicture("bloeberig.be")
	assert.Equal(t, data, picDataAfterSave)
}

func TestDelete(t *testing.T) {
	db := NewMentionRepo(conf)
	wm := mf.Mention{
		Target: "https://pussycat.com/coolpussy.html",
	}
	db.Save(wm, &mf.IndiewebData{
		Name: "lolz",
	})
	db.Delete(wm)

	results := db.GetAll("pussycat.com")
	assert.Equal(t, 0, len(results.Data))
}

func TestUpdateSince(t *testing.T) {
	db := NewMentionRepo(conf)
	nowStamp := time.Date(2020, 10, 13, 14, 15, 0, 0, time.UTC)

	db.UpdateSince("pussycat.com", nowStamp)
	since, err := db.Since("pussycat.com")

	assert.NoError(t, err)
	assert.Equal(t, nowStamp, since)
}

func TestSinceFirstTimeIsEmptytime(t *testing.T) {
	db := NewMentionRepo(conf)
	since, err := db.Since("pussycat.com")

	assert.Equal(t, buntdb.ErrNotFound, err)
	assert.Equal(t, time.Time{}, since)
}

func TestGet(t *testing.T) {
	db := NewMentionRepo(conf)
	wm := mf.Mention{
		Target: "https://pussycat.com/coolpussy.html",
	}
	db.Save(wm, &mf.IndiewebData{
		Name: "lolz",
	})

	result := db.Get(wm)
	assert.Equal(t, "lolz", result.Name)
}

func BenchmarkMentionRepoBunt_GetAll(b *testing.B) {
	defer os.Remove("test.db")
	db := NewMentionRepo(&common.Config{
		ConString: "test.db",
		AllowedWebmentionSources: []string{
			"pussycat.com",
		},
	})

	items := 10000
	fmt.Printf(" -- Saving %d items\n", items)
	for n := 0; n < items; n++ {
		db.Save(mf.Mention{
			Source: fmt.Sprintf("https://blahsource.com/%d/ding.html", n),
			Target: fmt.Sprintf("https://pussycat.com/%d/ding.html", n),
		}, &mf.IndiewebData{
			Name: fmt.Sprintf("benchmark %d", n),
			Author: mf.IndiewebAuthor{
				Name: fmt.Sprintf("author %d", n),
			},
		})
	}

	b.Run(fmt.Sprintf(" -- Benchmark Get All for #%d\n", b.N), func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			db.GetAll("pussycat.com")
		}
	})
}

func TestGetAllAndSaveSomeJson(t *testing.T) {
	db := NewMentionRepo(conf)
	db.Save(mf.Mention{
		Target: "https://pussycat.com/coolpussy.html",
	}, &mf.IndiewebData{
		Name: "lolz",
	})

	results := db.GetAll("pussycat.com")
	assert.Equal(t, 1, len(results.Data))
	assert.Equal(t, "lolz", results.Data[0].Name)
}

func TestGetFiltersBasedOnDomain(t *testing.T) {
	db := NewMentionRepo(conf)
	db.Save(mf.Mention{
		Target: "https://pussycat.com/coolpussy.html",
	}, &mf.IndiewebData{
		Name: "lolz",
	})
	db.Save(mf.Mention{
		Target: "https://dingeling.com/dogshateus.html",
	}, &mf.IndiewebData{
		Name: "amaigat",
	})

	results := db.GetAll("pussycat.com")
	assert.Equal(t, 1, len(results.Data))
	assert.Equal(t, "lolz", results.Data[0].Name)
}
