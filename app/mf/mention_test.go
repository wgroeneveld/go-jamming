package mf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDomainParseFromTarget(t *testing.T) {
	wm := Mention{
		Source: "source",
		Target: "http://patat.be/frietjes/zijn/lekker",
	}

	assert.Equal(t, "patat.be", wm.Domain())
}
