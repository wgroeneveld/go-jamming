package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
