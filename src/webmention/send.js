
const got = require('got')
const { collect } = require('./rsslinkcollector')
const { discover } = require('./../linkdiscoverer')
const { sendPingbackToEndpoint } = require('./../pingback/send')

async function sendWebmentionToEndpoint(endpoint, source, target) {
	await got.post(endpoint, {
		contentType: "x-www-form-urlencoded",
		form: {
			source,
			target
		},
		retry: {
			limit: 5,
			methods: ["POST"]
		}
	})
	console.log(` OK: webmention@${endpoint}, sent: source ${source}, target ${target}`)
}

async function mention(opts) {
	const { source, target } = opts
	const endpoint = await discover(target)
	const sendMention = {
		"webmention": sendWebmentionToEndpoint,
		"pingback": sendPingbackToEndpoint,
		"unknown": async function() {}
	}
	await sendMention[endpoint.type](endpoint.link, source, target)
}

async function parseRssFeed(xml, since) {
	const linksToMention = collect(xml, since)
		.map(el => el.hrefs
			// this strips relative URLs; could be a feature to also send these to own domain?
			.filter(href => href.startsWith('http'))
			.map(href => {
			return {
				// SOURCE is own domain this time, TARGET = outbound
				target: href,
				source: el.link
			}
		}))
		.flat()

	await Promise.all(linksToMention.map(mention))
}


async function send(domain, since) {
	const feed = await got(`https://${domain}/index.xml`)
	await parseRssFeed(feed.body, since)
}

module.exports = {
	send
}
