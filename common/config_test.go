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
	assert.Contains(t, config.AllowedWebmentionSources, "mycooldomain.com")
}

func TestReadFromJsonWithCorrectJsonData(t *testing.T) {
	confString := `{
		  "port": 1337,
		  "host": "localhost",
          "baseURL": "https://jam.brainbaking.com/",
		  "token": "miauwkes",
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

func TestWhitelist(t *testing.T) {
	conf := Config{
		Whitelist: []string{
			"youtube.com",
		},
		BaseURL:                  "https://jam.brainbaking.com/",
		Port:                     123,
		Token:                    "token",
		AllowedWebmentionSources: []string{"blah.com"},
	}
	t.Cleanup(func() {
		os.Remove("config.json")
	})

	conf.AddToWhitelist("dinges.be")
	assert.Contains(t, conf.Whitelist, "dinges.be")
	assert.Equal(t, 2, len(conf.Whitelist))

	confFromFile := Configure()
	assert.Contains(t, confFromFile.Whitelist, "dinges.be")
	assert.Equal(t, 2, len(confFromFile.Whitelist))
}

func TestAddToBlacklistNotYetAddsToListAndSaves(t *testing.T) {
	conf := Config{
		Blacklist: []string{
			"youtube.com",
		},
		BaseURL:                  "https://jam.brainbaking.com/",
		Port:                     123,
		Token:                    "token",
		AllowedWebmentionSources: []string{"blah.com"},
	}
	t.Cleanup(func() {
		os.Remove("config.json")
	})

	conf.AddToBlacklist("dinges.be")
	assert.Contains(t, conf.Blacklist, "dinges.be")
	assert.Equal(t, 2, len(conf.Blacklist))

	confFromFile := Configure()
	assert.Contains(t, confFromFile.Blacklist, "dinges.be")
	assert.Equal(t, 2, len(confFromFile.Blacklist))
}

func TestAddToBlacklistAlreadyAddedDoNotAddAgain(t *testing.T) {
	conf := Config{
		Blacklist: []string{
			"youtube.com",
		},
		Port:                     123,
		Token:                    "token",
		AllowedWebmentionSources: []string{"blah.com"},
	}
	t.Cleanup(func() {
		os.Remove("config.json")
	})

	conf.AddToBlacklist("youtube.com")
	assert.Contains(t, conf.Blacklist, "youtube.com")
	assert.Equal(t, 1, len(conf.Blacklist))
}

func TestIsWhitelisted(t *testing.T) {
	cases := []struct {
		label    string
		url      string
		expected bool
	}{
		{
			"do not whitelist if domain is part of relative url",
			"https://brainbaking.com/post/youtube.com-sucks",
			false,
		},
		{
			"whitelist if https domain is on the list",
			"https://youtube.com/stuff",
			true,
		},
		{
			"whitelist if http domain is on the list",
			"http://youtube.com/stuff",
			true,
		},
		{
			"do not whitelist if relative url",
			"/youtube.com",
			false,
		},
	}

	conf := Config{
		Whitelist: []string{
			"youtube.com",
		},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			assert.Equal(t, tc.expected, conf.IsWhitelisted(tc.url))
		})
	}
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
