package db

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"os"
)

type MentionRepo interface {
	// InModeration saves the mention data to the in moderation db to approve or reject later.
	// Returns the key or a marshal/persist error.
	InModeration(key mf.Mention, data *mf.IndiewebData) (string, error)
	// Save saves the mention to the approved db.
	// Returns the key or a marshal/persist error.
	Save(key mf.Mention, data *mf.IndiewebData) (string, error)
	// Delete removes a possibly present mention from the approved db by key.
	// Ignores but logs possible errors.
	Delete(key mf.Mention)
	// Approve saves the mention to the approved database and deletes the one in moderation.
	// If the key is invalid, it returns nil.
	Approve(key string) *mf.IndiewebData
	// Reject removes the in moderation key from the db and returns the deleted entry
	// If the key is invalid, it returns nil.
	Reject(key string) *mf.IndiewebData

	// Get returns a single unmarshalled json value based on the approved mention key in the db.
	// It returns the unmarshalled result or nil if something went wrong.
	Get(key mf.Mention) *mf.IndiewebData
	// GetAll returns a wrapped data result for all approved mentions for a particular domain.
	GetAll(domain string) mf.IndiewebDataResult
	// GetAll returns a wrapped data result for all to approve mentions for a particular domain.
	GetAllToModerate(domain string) mf.IndiewebDataResult

	// CleanupSpam removes potential blacklisted spam from the approved database by checking the url of each entry.
	CleanupSpam(domain string, blacklist []string)

	// SavePicture saves the picture byte data in the approved database and returns a key or error.
	SavePicture(bytes string, domain string) (string, error)
	// GetPicture returns a byte slice (or nil if unknown) from the approved database for a particular source domain.
	GetPicture(domain string) []byte
	// LastSentMention fetches the last known RSS link where mentions were sent from the approved db.
	// Returns an empty string if an error occured.
	LastSentMention(domain string) string
	// UpdateLastSentMention updates the last sent mention link in the approved db. Logs but ignores errors.
	UpdateLastSentMention(domain string, lastSent string)
}

type MentionRepoWrapper struct {
	toApproveRepo *mentionRepoBunt
	approvedRepo  *mentionRepoBunt
}

// Save saves the data to the
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

func (m MentionRepoWrapper) Approve(keyInModeration string) *mf.IndiewebData {
	toApprove := m.toApproveRepo.getByKey(keyInModeration)
	m.approvedRepo.saveByKey(keyInModeration, toApprove)
	m.toApproveRepo.deleteByKey(keyInModeration)
	return toApprove
}

func (m MentionRepoWrapper) Reject(keyInModeration string) *mf.IndiewebData {
	toReject := m.toApproveRepo.getByKey(keyInModeration)
	m.toApproveRepo.deleteByKey(keyInModeration)
	return toReject
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
