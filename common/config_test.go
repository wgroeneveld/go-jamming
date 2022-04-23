package common

import (
	"github.com/stretchr/testify/assert"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"
)

func cleanupConfig() {
	os.Remove("config.json")
}

func TestReadFromJsonMalformedReversToDefaults(t *testing.T) {
	ioutil.WriteFile("config.json", []byte("dinges"), fs.ModePerm)
	t.Cleanup(cleanupConfig)

	config := Configure()
	assert.Contains(t, config.AllowedWebmentionSources, "brainbaking.com")
}

func TestReadFromJsonWithCorrectJsonData(t *testing.T) {
	confString := `{
		  "port": 1337,
		  "host": "localhost",
		  "token": "miauwkes",
		  "utcOffset": 60,
		  "allowedWebmentionSources":  [
			"snoopy.be"
		  ],
		  "blacklist":  [
			"youtube.com"
		  ]
		}`
	ioutil.WriteFile("config.json", []byte(confString), fs.ModePerm)
	t.Cleanup(cleanupConfig)

	config := Configure()
	assert.Contains(t, config.AllowedWebmentionSources, "snoopy.be")
	assert.Equal(t, 1, len(config.AllowedWebmentionSources))
}

func TestSaveAfterAddingANewBlacklistEntry(t *testing.T) {
	t.Cleanup(cleanupConfig)

	config := Configure()
	config.AddToBlacklist("somethingnew.be")
	config.Save()

	newConfig := Configure()
	assert.Contains(t, newConfig.Blacklist, "somethingnew.be")
}

func TestAddToBlacklistNotYetAddsToList(t *testing.T) {
	conf := Config{
		Blacklist: []string{
			"youtube.com",
		},
	}

	conf.AddToBlacklist("dinges.be")
	assert.Contains(t, conf.Blacklist, "dinges.be")
	assert.Equal(t, 2, len(conf.Blacklist))
}

func TestAddToBlacklistAlreadyAddedDoNotAddAgain(t *testing.T) {
	conf := Config{
		Blacklist: []string{
			"youtube.com",
		},
	}

	conf.AddToBlacklist("youtube.com")
	assert.Contains(t, conf.Blacklist, "youtube.com")
	assert.Equal(t, 1, len(conf.Blacklist))
}

func TestIsBlacklisted(t *testing.T) {
	cases := []struct {
		label    string
		url      string
		expected bool
	}{
		{
			"do not blacklist if domain is part of relative url",
			"https://brainbaking.com/post/youtube.com-sucks",
			false,
		},
		{
			"blacklist if https domain is on the list",
			"https://youtube.com/stuff",
			true,
		},
		{
			"blacklist if http domain is on the list",
			"http://youtube.com/stuff",
			true,
		},
		{
			"do not blacklist if relative url",
			"/youtube.com",
			false,
		},
	}

	conf := Config{
		Blacklist: []string{
			"youtube.com",
		},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			assert.Equal(t, tc.expected, conf.IsBlacklisted(tc.url))
		})
	}
}
