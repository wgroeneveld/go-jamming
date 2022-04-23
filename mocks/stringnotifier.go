package mocks

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/notifier"
	"brainbaking.com/go-jamming/common"
)

type StringNotifier struct {
	Output string
	Conf   *common.Config
}

func (sn *StringNotifier) NotifyReceived(wm mf.Mention, indieweb *mf.IndiewebData) {
	sn.Output = notifier.BuildNotification(wm, indieweb, sn.Conf)
}
