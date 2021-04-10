package send

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/pingback/send"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/rest"
	"github.com/rs/zerolog/log"
	"time"
)

type Sender struct {
	RestClient rest.Client
	Conf       *common.Config
}

func (snder *Sender) Send(domain string, since string) {
	log.Info().Str("domain", domain).Str("since", since).Msg(` OK: someone wants to send mentions`)
	feed, err := snder.RestClient.GetBody("https://" + domain + "/index.xml")
	if err != nil {
		log.Err(err).Str("domain", domain).Msg("Unable to retrieve RSS feed, aborting send")
		return
	}

	snder.parseRssFeed(feed, common.IsoToTime(since))
}

func (snder *Sender) parseRssFeed(feed string, since time.Time) {

}

func mention() {
	pingbackSender := &send.Sender{
		RestClient: nil,
	}
	pingbackSender.SendPingbackToEndpoint("endpoint", mf.Mention{})
}
