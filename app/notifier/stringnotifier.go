package notifier

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
)

type StringNotifier struct {
	Output string
	Conf   *common.Config
}

func (sn *StringNotifier) NotifyInModeration(wm mf.Mention, data *mf.IndiewebData) error {
	sn.Output = "in moderation!"
	return nil
}
func (sn *StringNotifier) NotifyReceived(wm mf.Mention, data *mf.IndiewebData) error {
	sn.Output = "received!"
	return nil
}
