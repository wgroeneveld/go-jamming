package send

import (
	"github.com/rs/zerolog/log"
	"github.com/wgroeneveld/go-jamming/app/mf"
	"github.com/wgroeneveld/go-jamming/app/pingback/send"
	"github.com/wgroeneveld/go-jamming/common"
	"github.com/wgroeneveld/go-jamming/rest"
)

type Sender struct {
	RestClient rest.Client
	Conf       *common.Config
}

func (snder *Sender) Send(domain string, since string) {
	log.Info().Str("domain", domain).Str("since", since).Msg(` OK: someone wants to send mentions`)
}

func mention() {
	pingbackSender := &send.Sender{
		RestClient: nil,
	}
	pingbackSender.SendPingbackToEndpoint("endpoint", mf.Mention{})
}
