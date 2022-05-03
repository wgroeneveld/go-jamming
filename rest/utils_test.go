package rest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestIsRealImage(t *testing.T) {
	cases := []struct {
		label    string
		imgpath  string
		expected bool
	}{
		{
			"jpeg is a valid image",
			"../mocks/picture.jpg",
			true,
		},
		{
			"bmp is a valid image",
			"../mocks/picture.bmp",
			true,
		},
		{
			"xml is not a valid image",
			"../mocks/index.xml",
			false,
		},
		{
			"empty data is not a valid image",
			"",
			false,
		},
		{
			"png is a valid image",
			"../mocks/picture.png",
			true,
		},
		{
			"gif is a valid image",
			"../mocks/picture.gif",
			true,
		},
		{
			"webp is a valid image",
			"../mocks/picture.webp",
			true,
		},
		{
			"tiff is a valid image",
			"../mocks/picture.tiff",
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			data, _ := ioutil.ReadFile(tc.imgpath)
			fmt.Printf("Path: %s, Data: % x\n", tc.imgpath, data)

			assert.Equal(t, tc.expected, IsRealImage(data))
		})
	}
}

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
		{
			"parse from localhost domain without extension",
			"https://localhost:1313/stuff",
			"localhost",
		},
		{
			"UK domain with two dots after the name",
			"https://minutestomidnight.co.uk/blog/article.html",
			"minutestomidnight.co.uk",
		},
		{
			"UK domain with subdomain",
			"https://www.minutestomidnight.co.uk/blog/article.html",
			"minutestomidnight.co.uk",
		},
		{
			"malformed http string with too little slashes simply returns same URL",
			"https:*groovy.bla/stuff",
			"https:*groovy.bla/stuff",
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			assert.Equal(t, tc.expected, Domain(tc.target))
		})
	}
}
