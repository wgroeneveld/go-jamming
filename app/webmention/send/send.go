package send

import (
	"github.com/wgroeneveld/go-jamming/app/mf"
	"github.com/wgroeneveld/go-jamming/app/pingback/send"
	"github.com/wgroeneveld/go-jamming/common"
	"github.com/wgroeneveld/go-jamming/rest"
)

type Sender struct {
	RestClient rest.Client
	Conf       *common.Config
}

func mention() {
	pingbackSender := &send.Sender{
		RestClient: nil,
	}
	pingbackSender.SendPingbackToEndpoint("endpoint", mf.Mention{})
}
