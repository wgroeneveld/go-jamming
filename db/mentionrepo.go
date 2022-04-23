package db

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"os"
)

type MentionRepo interface {
	InModeration(key mf.Mention, data *mf.IndiewebData) (string, error)
	Save(key mf.Mention, data *mf.IndiewebData) (string, error)
	Delete(key mf.Mention)
	Approve(key mf.Mention)
	Reject(key mf.Mention)

	Get(key mf.Mention) *mf.IndiewebData
	GetAll(domain string) mf.IndiewebDataResult
	GetAllToModerate(domain string) mf.IndiewebDataResult

	CleanupSpam(domain string, blacklist []string)

	SavePicture(bytes string, domain string) (string, error)
	GetPicture(domain string) []byte
	LastSentMention(domain string) string
	UpdateLastSentMention(domain string, lastSent string)
}

type MentionRepoWrapper struct {
	toApproveRepo *mentionRepoBunt
	approvedRepo  *mentionRepoBunt
}

func (m MentionRepoWrapper) Save(key mf.Mention, data *mf.IndiewebData) (string, error) {
	return m.approvedRepo.Save(key, data)
}

func (m MentionRepoWrapper) InModeration(key mf.Mention, data *mf.IndiewebData) (string, error) {
	return m.toApproveRepo.Save(key, data)
}

func (m MentionRepoWrapper) SavePicture(bytes string, domain string) (string, error) {
	return m.approvedRepo.SavePicture(bytes, domain)
}

func (m MentionRepoWrapper) Delete(key mf.Mention) {
	m.approvedRepo.Delete(key)
}

func (m MentionRepoWrapper) Approve(keyInModeration mf.Mention) {
	toApprove := m.toApproveRepo.Get(keyInModeration)
	m.Save(keyInModeration, toApprove)
	m.toApproveRepo.Delete(keyInModeration)
}

func (m MentionRepoWrapper) Reject(keyInModeration mf.Mention) {
	m.toApproveRepo.Delete(keyInModeration)
}

func (m MentionRepoWrapper) CleanupSpam(domain string, blacklist []string) {
	m.approvedRepo.CleanupSpam(domain, blacklist)
}

func (m MentionRepoWrapper) LastSentMention(domain string) string {
	return m.approvedRepo.LastSentMention(domain)
}

func (m MentionRepoWrapper) UpdateLastSentMention(domain string, lastSent string) {
	m.approvedRepo.UpdateLastSentMention(domain, lastSent)
}

func (m MentionRepoWrapper) Get(key mf.Mention) *mf.IndiewebData {
	return m.approvedRepo.Get(key)
}

func (m MentionRepoWrapper) GetPicture(domain string) []byte {
	return m.approvedRepo.GetPicture(domain)
}

func (m MentionRepoWrapper) GetAll(domain string) mf.IndiewebDataResult {
	return m.approvedRepo.GetAll(domain)
}

func (m MentionRepoWrapper) GetAllToModerate(domain string) mf.IndiewebDataResult {
	return m.toApproveRepo.GetAll(domain)
}

// NewMentionRepo returns a wrapper to two different mentionRepoBunt instances
// Depending on the to approve or approved mention, it will be saved in another file.
func NewMentionRepo(c *common.Config) *MentionRepoWrapper {
	return &MentionRepoWrapper{
		toApproveRepo: newMentionRepoBunt("mentions_toapprove.db", c.AllowedWebmentionSources),
		approvedRepo:  newMentionRepoBunt("mentions.db", c.AllowedWebmentionSources),
	}
}

// Purge removes all database files from disk.
// This is dangerous in production and should be used as a shorthand in tests!
func Purge() {
	os.Remove("mentions_toapprove.db")
	os.Remove("mentions.db")
}
