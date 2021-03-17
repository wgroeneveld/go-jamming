
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
	const links = description.match(/href="([^"]*")/g)
		.map(match => match.replace("href=", "").replace(/\"/g, ""))
		.filter(match => !(/\.(gif|zip|rar|bz2|gz|7z|jpe?g|tiff?|png|webp|bmp)$/i).test(match))
		.filter(match => !config.disallowedWebmentionDomains.some(domain => match.indexOf(domain) >= 0))
	return [...new Set(links)]
}

/**
* a typical RSS item looks like this:
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
      '            By <a href="/about">Wouter Groeneveld</a> on 16 March 2021.\n' +
      '          </p>\n' +
      '          '
  }
**/ 
function collect(xml, since = '') {
  const root = parser.parse(xml, parseOpts).rss.channel
  const sinceDate = dayjs(since)

  // example pubDate format: Tue, 16 Mar 2021 17:07:14 +0000
  const sincePubDate = (date) => {
  	if(!sinceDate.isValid()) return true
  	const pubDate = dayjs(date.split(", ")[1], "DD MMM YYYY HH:mm:ss ZZ")
  	if(!pubDate.isValid()) return true
  	return sinceDate < pubDate
  }

  const entries = root.item.filter ? root.item : [root.item]

  return entries
  	.filter(item => sincePubDate(item.pubDate))
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
