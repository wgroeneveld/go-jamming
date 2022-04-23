// A database wrapper package for BuntDB that persists indieweb (meta)data.
// Most functions silently suppress errors as with consistent states, it would be impossible.
package db

import (
	"brainbaking.com/go-jamming/app/mf"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/buntdb"
	"strings"
)

type mentionRepoBunt struct {
	db *buntdb.DB
}

// CleanupSpam removes potential blacklisted spam from the webmention database by checking the url of each entry.
func (r *mentionRepoBunt) CleanupSpam(domain string, blacklist []string) {
	for _, mention := range r.GetAll(domain).Data {
		for _, blacklisted := range blacklist {
			if strings.Contains(mention.Url, blacklisted) {
				r.Delete(mention.AsMention())
			}
		}
	}
}

// UpdateLastSentMention updates the last sent mention link. Logs but ignores errors.
func (r *mentionRepoBunt) UpdateLastSentMention(domain string, lastSentMentionLink string) {
	err := r.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(lastSentKey(domain), lastSentMentionLink, nil)
		return err
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateLastSentMention: unable to save")
	}
}

// LastSentMention fetches the last known RSS link where mentions were sent, or an empty string if an error occured.
func (r *mentionRepoBunt) LastSentMention(domain string) string {
	var lastSent string
	err := r.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(lastSentKey(domain))
		lastSent = val
		return err
	})
	if err != nil {
		log.Error().Err(err).Msg("LastSentMention: unable to retrieve last sent, reverting to empty")
		return ""
	}
	return lastSent
}

func lastSentKey(domain string) string {
	return fmt.Sprintf("%s:lastsent", domain)
}

// Delete removes a possibly present mention by key. Ignores but logs possible errors.
func (r *mentionRepoBunt) Delete(wm mf.Mention) {
	key := r.mentionToKey(wm)
	err := r.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})
	if err != nil {
		log.Warn().Err(err).Str("key", key).Stringer("wm", wm).Msg("Unable to delete")
	} else {
		log.Debug().Str("key", key).Stringer("wm", wm).Msg("Deleted.")
	}
}

func (r *mentionRepoBunt) SavePicture(bytes string, domain string) (string, error) {
	key := pictureKey(domain)
	err := r.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, bytes, nil)
		return err
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func pictureKey(domain string) string {
	return fmt.Sprintf("%s:picture", domain)
}

// Save saves the mention by marshalling data. Returns the key or a marshal/persist error.
func (r *mentionRepoBunt) Save(wm mf.Mention, data *mf.IndiewebData) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	key := r.mentionToKey(wm)
	err = r.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, string(jsonData), nil)
		return err
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (r *mentionRepoBunt) mentionToKey(wm mf.Mention) string {
	return fmt.Sprintf("%s:%s", wm.Key(), wm.Domain())
}

// Get returns a single unmarshalled json value based on the mention key.
// It returns the unmarshalled result or nil if something went wrong.
func (r *mentionRepoBunt) Get(wm mf.Mention) *mf.IndiewebData {
	var data mf.IndiewebData
	key := r.mentionToKey(wm)
	err := r.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(val), &data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("repo get: unable to retrieve key")
		return nil
	}
	return &data
}

func (r *mentionRepoBunt) GetPicture(domain string) []byte {
	var data []byte
	key := pictureKey(domain)
	err := r.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		data = []byte(val)
		return nil
	})
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("repo getpicture: unable to retrieve key")
		return nil
	}
	return data
}

// GetAll returns a wrapped data result for all mentions for a particular domain.
// Intentionally ignores marshal errors, db should be consistent!
// Warning, this will potentially marshall 10k strings! See benchmark test.
func (r *mentionRepoBunt) GetAll(domain string) mf.IndiewebDataResult {
	var data []*mf.IndiewebData
	err := r.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(domain, func(key, value string) bool {
			var result mf.IndiewebData
			json.Unmarshal([]byte(value), &result)
			data = append(data, &result)
			return true
		})
	})

	if err != nil {
		log.Error().Err(err).Msg("get all: failed to ascend from view")
		return mf.ResultFailure(data)
	}
	return mf.ResultSuccess(data)
}

// NewMentionRepoBunt opens a database connection using default buntdb settings.
// It also creates necessary indexes based on the passed domain config.
// This panics if it cannot open the db.
func newMentionRepoBunt(conString string, allowedWebmentionSources []string) *mentionRepoBunt {
	approvedRepo := &mentionRepoBunt{}
	db, err := buntdb.Open(conString)
	if err != nil {
		log.Fatal().Str("constr", conString).Msg("new mention repo: cannot open db")
	}
	approvedRepo.db = db

	for _, domain := range allowedWebmentionSources {
		db.CreateIndex(domain, fmt.Sprintf("*:%s", domain), buntdb.IndexString)
	}

	return approvedRepo
}
