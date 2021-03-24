
async function send(domain, since) {
	const feed = await got(`https://${domain}/index.xml`)
	await parseRssFeed(feed.body, since)
}

module.exports = {
	send
}
