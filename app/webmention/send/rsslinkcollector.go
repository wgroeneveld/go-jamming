package send

import (
	"brainbaking.com/go-jamming/app/rss"
	"brainbaking.com/go-jamming/common"
	"regexp"
	"time"
)

type RSSItem struct {
	link  string
	hrefs []string
}

/**
* a typical RSS item looks like this:
-- if <time/> found in body, assume it's a lastmod update timestamp!
 {
    title: '@celia @kev I have read both you and Kev&#39;s post on...',
    link: 'https://brainbaking.com/notes/2021/03/16h17m07s14/',
    comments: 'https://brainbaking.com/notes/2021/03/16h17m07s14/#commento',
    pubDate: 'Tue, 16 Mar 2021 17:07:14 +0000',
    author: 'Wouter Groeneveld',
    guid: {
      '#text': 'https://brainbaking.com/notes/2021/03/16h17m07s14/',
      '@_isPermaLink': 'true'
    },
    description: ' \n' +
      '          \n' +
      '\n' +
      '          <p><span class="h-card"><a class="u-url mention" data-user="A5GVjIHI6MH82H6iLQ" href="https://fosstodon.org/@celia" rel="ugc">@<span>celia</span></a></span> <span class="h-card"><a class="u-url mention" data-user="A54b8g0RBaIgjzczMu" href="https://fosstodon.org/@kev" rel="ugc">@<span>kev</span></a></span> I have read both you and Kev&rsquo;s post on this and agree on some points indeed! But I&rsquo;m not yet ready to give up webmentions. As an academic, the idea of citing/mentioning each other is very alluring 🤓. Plus, I needed an excuse to fiddle some more with JS&hellip; <br><br>As much as I loved using Wordpress before, I can&rsquo;t imagine going back to writing stuff in there instead of in markdown. Gotta keep the workflow short, though. Hope it helps you focus on what matters - content!</p>\n' +
      '\n' +
      '\n' +
      '          <p>\n' +
      '            By <a href="/about">Wouter Groeneveld</a> on <time datetime='2021-03-20'>20 March 2021</time>.\n' +
      '          </p>\n' +
      '          '
  }
**/
func (snder *Sender) Collect(xml string, since time.Time) ([]RSSItem, error) {
	feed, err := rss.ParseFeed([]byte(xml))
	if err != nil {
		return nil, err
	}
	var items []RSSItem
	for _, rssitem := range feed.ItemList {
		if since.IsZero() || since.Before(rssitem.PubDateAsTime()) {
			items = append(items, RSSItem{
				link:  rssitem.Link,
				hrefs: snder.collectUniqueHrefsFromDescription(rssitem.Description),
			})
		}
	}
	return items, nil
}

func (snder *Sender) collectUniqueHrefsFromDescription(html string) []string {
	r := regexp.MustCompile(`href="(.+?)"`)
	ext := regexp.MustCompile(`\.(gif|zip|rar|bz2|gz|7z|jpe?g|tiff?|png|webp|bmp)$`)
	urlmap := common.NewSet()

	for _, match := range r.FindAllStringSubmatch(html, -1) {
		url := match[1] // [0] is the match of the entire expression, [1] is the capture group
		if !ext.MatchString(url) && !snder.Conf.ContainsDisallowedDomain(url) {
			urlmap.Add(url)
		}
	}

	return urlmap.Keys()
}
