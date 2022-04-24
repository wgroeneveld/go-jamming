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
	SourceDomain string
	Source       string
	Content      string
	Target       string
	AdminURL     string
	ApproveURL   string
	RejectURL    string
}

type Notifier interface {
	NotifyReceived(wm mf.Mention, data *mf.IndiewebData)
}

// BuildNotification returns a HTML (string template) representation of the Mention to notify the admin.
func BuildNotification(wm mf.Mention, data *mf.IndiewebData, cnf *common.Config) string {
	acceptUrl := fmt.Sprintf("%sadmin/approve/%s/%s", cnf.BaseURL, cnf.Token, wm.Key())
	rejectUrl := fmt.Sprintf("%sadmin/reject/%s/%s", cnf.BaseURL, cnf.Token, wm.Key())
	adminUrl := fmt.Sprintf("%sadmin/%s", cnf.BaseURL, cnf.Token)

	var buff bytes.Buffer
	notificationTmpl.Execute(&buff, notificationData{
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
