package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncludes(t *testing.T) {
	cases := []struct {
		label     string
		arr       []string
		searchstr string
		expected  bool
	}{
		{
			"element in array",
			[]string{"one", "two"},
			"two",
			true,
		},
		{
			"element not in array",
			[]string{"one", "two"},
			"three",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			assert.Equal(t, tc.expected, Includes(tc.arr, tc.searchstr))
		})
	}
}
