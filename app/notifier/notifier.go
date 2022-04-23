package notifier

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"fmt"
)

type Notifier interface {
	NotifyReceived(wm mf.Mention, data *mf.IndiewebData)
}

// BuildNotification returns a string representation of the Mention to notify the admin.
func BuildNotification(wm mf.Mention, data *mf.IndiewebData, cnf *common.Config) string {
	enter := "\n"
	acceptUrl := fmt.Sprintf("%sadmin/approve/%s/%s", cnf.BaseURL, cnf.Token, wm.Key())
	rejectUrl := fmt.Sprintf("%sadmin/reject/%s/%s", cnf.BaseURL, cnf.Token, wm.Key())

	return fmt.Sprintf("Hi admin, %s%s,A webmention was received: %sSource %s, Target %s%sContent: %s%s%sAccept? %s%sReject? %s%sCheerio, your go-jammin' thing.",
		enter, enter, enter,
		wm.Source, wm.Target, enter,
		data.Content, enter, enter,
		acceptUrl, enter, rejectUrl, enter)
}
