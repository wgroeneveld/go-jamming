package send

import (
	"github.com/rs/zerolog/log"
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/pingback/send"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/rest"
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
