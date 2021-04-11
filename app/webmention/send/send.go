package send

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/pingback/send"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/rest"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Sender struct {
	RestClient rest.Client
	Conf       *common.Config
}

func (snder *Sender) Send(domain string, since string) {
	log.Info().Str("domain", domain).Str("since", since).Msg(` OK: someone wants to send mentions`)
	_, feed, err := snder.RestClient.GetBody("https://" + domain + "/index.xml")
	if err != nil {
		log.Err(err).Str("domain", domain).Msg("Unable to retrieve RSS feed, aborting send")
		return
	}

	snder.parseRssFeed(feed, common.IsoToTime(since))
	log.Info().Str("domain", domain).Str("since", since).Msg(` OK: sending done.`)
}

func (snder *Sender) parseRssFeed(feed string, since time.Time) {
	items, err := snder.Collect(feed, since)
	if err != nil {
		log.Err(err).Msg("Unable to parse RSS fed, aborting send")
		return
	}

	var wg sync.WaitGroup
	for _, item := range items {
		for _, href := range item.hrefs {
			mention := mf.Mention{
				// SOURCE is own domain this time, TARGET = outbound
				Source: item.link,
				Target: href,
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				snder.sendMention(mention)
			}()
		}
	}
	wg.Wait()
}

var mentionFuncs = map[string]func(snder *Sender, mention mf.Mention, endpoint string){
	TypeUnknown:    func(snder *Sender, mention mf.Mention, endpoint string) {},
	TypeWebmention: sendMentionAsWebmention,
	TypePingback:   sendMentionAsPingback,
}

func (snder *Sender) sendMention(mention mf.Mention) {
	endpoint, mentionType := snder.discover(mention.Target)
	mentionFuncs[mentionType](snder, mention, endpoint)
}

func sendMentionAsWebmention(snder *Sender, mention mf.Mention, endpoint string) {
	err := snder.RestClient.PostForm(endpoint, mention.AsFormValues())
	if err != nil {
		log.Err(err).Str("endpoint", endpoint).Str("wm", mention.String()).Msg("Webmention send failed")
		return
	}
	log.Info().Str("endpoint", endpoint).Str("wm", mention.String()).Msg("OK: webmention sent.")
}

func sendMentionAsPingback(snder *Sender, mention mf.Mention, endpoint string) {
	pingbackSender := &send.Sender{
		RestClient: snder.RestClient,
	}
	pingbackSender.SendPingbackToEndpoint(endpoint, mention)
}
