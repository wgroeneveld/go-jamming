package send

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/app/pingback/send"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/db"
	"brainbaking.com/go-jamming/rest"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
)

type Sender struct {
	RestClient rest.Client
	Conf       *common.Config
	Repo       db.MentionRepo
}

// SendSingle sends out webmentions serially for a single source.
// It does validate the relative path against the domain, which is supposed to be served using https.
func (snder *Sender) SendSingle(domain string, relSource string) {
	source := fmt.Sprintf("https://%s/%s", domain, relSource)
	log.Info().Str("url", source).Msg(` OK: someone wants to send a single mention`)

	_, html, err := snder.RestClient.GetBody(source)
	if err != nil {
		log.Err(err).Str("url", source).Msg("Unable to validate source, send aborted")
		return
	}

	for _, href := range snder.collectUniqueHrefsFromHtml(html) {
		if strings.HasPrefix(href, "http") {
			snder.sendMention(mf.Mention{
				Source: source,
				Target: href,
			})
		}
	}
}

// Send sends out multiple webmentions based on since and what's posted in the RSS feed.
// It first GETs domain/index.xml and goes from there.
func (snder *Sender) Send(domain string) {
	lastSent := snder.Repo.LastSentMention(domain)
	feedUrl := "https://" + domain + "/index.xml"

	log.Info().Str("domain", domain).Str("lastsent", lastSent).Msg(` OK: someone wants to send mentions`)
	_, feed, err := snder.RestClient.GetBody(feedUrl)
	if err != nil {
		log.Err(err).Str("url", feedUrl).Msg("Unable to retrieve RSS feed, send aborted")
		return
	}

	lastSent, err = snder.parseRssFeed(feed, lastSent)
	if err != nil {
		log.Err(err).Str("url", feedUrl).Msg("Unable to parse RSS feed, send aborted")
		return
	}

	snder.Repo.UpdateLastSentMention(domain, lastSent)
	log.Info().Str("domain", domain).Str("lastsent", lastSent).Msg(` OK: send processed.`)
}

func (snder *Sender) parseRssFeed(feed string, lastSentLink string) (string, error) {
	items, err := snder.Collect(feed, lastSentLink)
	if err != nil {
		return lastSentLink, err
	}
	if len(items) == 0 {
		return lastSentLink, nil
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
	// first item is the most recent one!
	return items[0].link, nil
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
