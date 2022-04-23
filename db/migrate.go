package db

import (
	"brainbaking.com/go-jamming/common"
)

// Migrate self-checks and executes necessary DB migrations, if any.
func Migrate() {
	cnf := common.Configure()
	repo := NewMentionRepo(cnf)

	// no migrations needed anymore/yet
	repo.approvedRepo.db.Shrink()
	repo.toApproveRepo.db.Shrink()
}
