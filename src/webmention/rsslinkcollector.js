
const parser = require("fast-xml-parser")
const config = require('./../config')

const dayjs = require('dayjs')
const customParseFormat = require('dayjs/plugin/customParseFormat')
dayjs.extend(customParseFormat)

const parseOpts = {
    ignoreAttributes: false
}

function collectHrefsFromDescription(description) {
	// first thought: use parser.parse() and traverse recursively. turned out to be way too slow.
	const linksMatch = description.match(/href="([^"]*")/g)
  if(!linksMatch) return []

  const links = linksMatch
		.map(match => match.replace("href=", "").replace(/\"/g, ""))
		.filter(match => !(/\.(gif|zip|rar|bz2|gz|7z|jpe?g|tiff?|png|webp|bmp)$/i).test(match))
		.filter(match => !config.disallowedWebmentionDomains.some(domain => match.indexOf(domain) >= 0))
	return [...new Set(links)]
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
function collect(xml, since = '') {
  const root = parser.parse(xml, parseOpts).rss.channel
  const sinceDate = dayjs(since)

  const enrichWithDateProperties = (item) => {
    // example pubDate format: Tue, 16 Mar 2021 17:07:14 +0000
    const rawpub = item.pubDate?.split(", ")?.[1]
    item.pubDate = rawpub ? dayjs(rawpub, "DD MMM YYYY HH:mm:ss ZZ") : dayjs()
    if(!item.pubDate.isValid()) item.pubDate = dayjs()

    const dateTimeMatch = item.description.match(/datetime="([^"]*")/g)
    // Selecting the first - could be dangerous. Living on the edge. Don't care. etc. 
    const rawlastmod = dateTimeMatch?.[0]?.replace("datetime=", "")?.replace(/\"/g, "")
    item.lastmodDate = rawlastmod ? dayjs(rawlastmod) : dayjs(0)

    return item
  }

  const sincePublicationDate = (item) => {
  	if(!sinceDate.isValid()) return true
  	
    return sinceDate < (item.lastmodDate > item.pubDate ? item.lastmodDate : item.pubDate)
  }

  const entries = root.item.filter ? root.item : [root.item]

  return entries
    .map(enrichWithDateProperties)
  	.filter(sincePublicationDate)
  	.map(item => {
  	return {
  		link: item.link,
  		hrefs: collectHrefsFromDescription(item.description)
  	}
  })
}

module.exports = {
	collect
}
