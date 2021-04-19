package rest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDomainParseFromTarget(t *testing.T) {
	cases := []struct {
		label    string
		target   string
		expected string
	}{
		{
			"parse from default http domain",
			"http://patat.be/frietjes/zijn/lekker",
			"patat.be",
		},
		{
			"parse from default https domain",
			"https://frit.be/patatjes/zijn/lekker",
			"frit.be",
		},
		{
			"parse from default https domain with www subdomain",
			"https://www.frit.be/patatjes/zijn/lekker",
			"frit.be",
		},
		{
			"parse from default https domain with some random subdomain",
			"https://mayonaise.frit.be/patatjes/zijn/lekker",
			"frit.be",
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			assert.Equal(t, tc.expected, Domain(tc.target))
		})
	}
}
