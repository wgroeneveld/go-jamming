package send

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/pingback/send"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/rest"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/fs"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

type Sender struct {
	RestClient rest.Client
	Conf       *common.Config
}

func (snder *Sender) sinceForDomain(domain string, since string) time.Time {
	if since != "" {
		return common.IsoToTime(since)
	}
	sinceConf, err := ioutil.ReadFile(fmt.Sprintf("%s/%s-since.txt", snder.Conf.DataPath, domain))
	if err != nil {
		log.Warn().Str("since", since).Msg("No query param, and no config found. Reverting to beginning of time...")
		return time.Time{}
	}
	return common.IsoToTime(string(sinceConf))
}

func (snder *Sender) saveSinceForDomain(domain string, since time.Time) {
	ioutil.WriteFile(fmt.Sprintf("%s/%s-since.txt", snder.Conf.DataPath, domain), []byte(common.TimeToIso(since)), fs.ModePerm)
}

func (snder *Sender) Send(domain string, since string) {
	snder.Conf.Lock(domain)
	defer snder.Conf.Unlock(domain)
	timeSince := snder.sinceForDomain(domain, since)
	log.Info().Str("domain", domain).Time("since", timeSince).Msg(` OK: someone wants to send mentions`)
	feedUrl := "https://" + domain + "/index.xml"
	_, feed, err := snder.RestClient.GetBody(feedUrl)
	if err != nil {
		log.Err(err).Str("url", feedUrl).Msg("Unable to retrieve RSS feed, send aborted")
		return
	}

	if err = snder.parseRssFeed(feed, timeSince); err != nil {
		log.Err(err).Str("url", feedUrl).Msg("Unable to parse RSS feed, send aborted")
		return
	}

	snder.saveSinceForDomain(domain, timeSince)
}

func (snder *Sender) parseRssFeed(feed string, since time.Time) error {
	items, err := snder.Collect(feed, since)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	sema := make(chan struct{}, 20)

	for _, item := range items {
		for _, href := range item.hrefs {
			if strings.HasPrefix(href, "http") {
				mention := mf.Mention{
					// SOURCE is own domain this time, TARGET = outbound
					Source: item.link,
					Target: href,
				}

				wg.Add(1)
				go func() {
					sema <- struct{}{}
					defer func() { <-sema }()
					defer wg.Done()
					snder.sendMention(mention)
				}()
			}
		}
	}
	wg.Wait()
	return nil
}

var mentionFuncs = map[string]func(snder *Sender, mention mf.Mention, endpoint string){
	typeUnknown:    func(snder *Sender, mention mf.Mention, endpoint string) {},
	typeWebmention: sendMentionAsWebmention,
	typePingback:   sendMentionAsPingback,
}

func (snder *Sender) sendMention(mention mf.Mention) {
	endpoint, mentionType := snder.discover(mention.Target)
	mentionFuncs[mentionType](snder, mention, endpoint)
}

func sendMentionAsWebmention(snder *Sender, mention mf.Mention, endpoint string) {
	err := snder.RestClient.PostForm(endpoint, mention.AsFormValues())
	if err != nil {
		log.Err(err).Str("endpoint", endpoint).Stringer("wm", mention).Msg("Webmention send failed")
		return
	}
	log.Info().Str("endpoint", endpoint).Stringer("wm", mention).Msg("OK: webmention sent.")
}

func sendMentionAsPingback(snder *Sender, mention mf.Mention, endpoint string) {
	pingbackSender := &send.Sender{
		RestClient: snder.RestClient,
	}
	pingbackSender.SendPingbackToEndpoint(endpoint, mention)
}
