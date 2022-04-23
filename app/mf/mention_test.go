package mf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTargetDomainDomain(t *testing.T) {
	wm := Mention{
		Source: "source",
		Target: "http://patat.be/frietjes/zijn/lekker",
	}

	assert.Equal(t, "patat.be", wm.TargetDomain())
}

func TestSourceDomain(t *testing.T) {
	wm := Mention{
		Source: "http://patat.be/frietjes/zijn/lekker",
		Target: "source",
	}

	assert.Equal(t, "patat.be", wm.SourceDomain())

}
