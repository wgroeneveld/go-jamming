package notifier

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"text/template"
)

import _ "embed"

//go:embed notification.html
var notificationTmplBytes []byte
var notificationTmpl *template.Template

func init() {
	var err error
	notificationTmpl, err = template.New("notification").Parse(string(notificationTmplBytes))
	if err != nil {
		log.Fatal().Err(err).Str("name", "notification").Msg("Template invalid")
	}
}

type notificationData struct {
	Action       string
	SourceDomain string
	Source       string
	Content      string
	Target       string
	AdminURL     string
	ApproveURL   string
	RejectURL    string
}

type Notifier interface {
	NotifyInModeration(wm mf.Mention, data *mf.IndiewebData) error
	NotifyReceived(wm mf.Mention, data *mf.IndiewebData) error
}

// buildReceivedMsg returns a HTML (string template) representation of the approved mention to notify the admin.
func buildReceivedMsg(wm mf.Mention, data *mf.IndiewebData, cnf *common.Config) string {
	adminUrl := adminUrl(cnf)
	var buff bytes.Buffer
	notificationTmpl.Execute(&buff, notificationData{
		Action:       "approved",
		Source:       wm.Source,
		Target:       wm.Target,
		Content:      data.Content,
		SourceDomain: wm.SourceDomain(),
		AdminURL:     adminUrl,
	})
	return buff.String()
}

// buildInModerationMsg returns a HTML (string template) representation of the in moderation mention to notify the admin.
func buildInModerationMsg(wm mf.Mention, data *mf.IndiewebData, cnf *common.Config) string {
	acceptUrl := acceptUrl(wm, cnf)
	rejectUrl := rejectUrl(wm, cnf)
	adminUrl := adminUrl(cnf)

	var buff bytes.Buffer
	notificationTmpl.Execute(&buff, notificationData{
		Action:       "in moderation",
		Source:       wm.Source,
		Target:       wm.Target,
		Content:      data.Content,
		SourceDomain: wm.SourceDomain(),
		ApproveURL:   acceptUrl,
		RejectURL:    rejectUrl,
		AdminURL:     adminUrl,
	})
	return buff.String()
}

func rejectUrl(wm mf.Mention, cnf *common.Config) string {
	return fmt.Sprintf("%sadmin/reject/%s/%s", cnf.BaseURL, cnf.Token, wm.Key())
}

func acceptUrl(wm mf.Mention, cnf *common.Config) string {
	return fmt.Sprintf("%sadmin/approve/%s/%s", cnf.BaseURL, cnf.Token, wm.Key())
}

func adminUrl(cnf *common.Config) string {
	return fmt.Sprintf("%sadmin/%s", cnf.BaseURL, cnf.Token)
}
