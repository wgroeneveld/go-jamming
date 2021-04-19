package db

import "brainbaking.com/go-jamming/common"

// Migrate self-checks and executes necessary DB migrations, if any.
func Migrate() {
	cnf := common.Configure()
	repo := NewMentionRepo(cnf)

	MigrateDataFiles(cnf, repo)
	MigratePictures(cnf, repo)
	repo.db.Shrink()
}
