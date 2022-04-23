package db

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	repoCnf = &common.Config{
		AllowedWebmentionSources: []string{
			"brainbaking.com",
		},
	}
)

func TestRejectUnknownKeyReturnsNil(t *testing.T) {
	repo := NewMentionRepo(repoCnf)
	t.Cleanup(Purge)

	result := repo.Reject("fuckballz")
	assert.Nil(t, result)
}

func TestApproveUnknownKeyReturnsNil(t *testing.T) {
	repo := NewMentionRepo(repoCnf)
	t.Cleanup(Purge)

	result := repo.Approve("fuckballz")
	assert.Nil(t, result)
}

func TestApproveAndRejectCases(t *testing.T) {
	cases := []struct {
		label                  string
		approve                bool
		expectedInModerationDb int
		expectedInMentionDb    int
	}{
		{
			"approve moves from the to moderate db to the mention db",
			true,
			0,
			1,
		},
		{
			"reject deletes from to moderate db and leaves mention db alone",
			false,
			0,
			0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			repo := NewMentionRepo(repoCnf)
			defer Purge()

			wm := mf.Mention{
				Source: "https://jefklakscodex.com/dinges",
				Target: "https://brainbaking.com/sjiekedinges.html",
			}
			data := &mf.IndiewebData{
				Name: "lolz",
			}
			repo.InModeration(wm, data)

			if tc.approve {
				repo.Approve(wm.Key())
			} else {
				repo.Reject(wm.Key())
			}

			allWms := repo.GetAll("brainbaking.com")
			allWmsToModerate := repo.GetAllToModerate("brainbaking.com")
			assert.Equal(t, tc.expectedInMentionDb, len(allWms.Data), "mention db expectation failed")
			assert.Equal(t, tc.expectedInModerationDb, len(allWmsToModerate.Data), "in moderation db expectation failed")
		})
	}
}
