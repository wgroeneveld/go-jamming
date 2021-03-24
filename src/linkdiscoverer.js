const got = require('got')
const { mf2 } = require("microformats-parser");

const baseUrlOf = (url) => {
	if(url.match(/\//g).length <= 2) {
		return url
	}
	const split = url.split('/')
	return split[0] + '//' + split[2]
}

const buildWebmentionHeaderLink = (link) => {
	// e.g. Link: <http://aaronpk.example/webmention-endpoint>; rel="webmention"
	return link
		.split(";")[0]
		.replace("<" ,"")
		.replace(">", "")
}

// see https://www.w3.org/TR/webmention/#sender-discovers-receiver-webmention-endpoint
async function discover(target) {
	try {
		const endpoint = await got(target)
		if(endpoint.headers.link?.indexOf("webmention") >= 0) {
			return {
				link: buildWebmentionHeaderLink(endpoint.headers.link),
				type: "webmention"
			}
		} else if(endpoint.headers["X-Pingback"]) {
			return {
				link: endpoint.headers["X-Pingback"],
				type: "pingback"
			}
		}

		const format = mf2(endpoint.body, {
			// this also complies with w3.org regulations: relative endpoint could be possible
			baseUrl: baseUrlOf(target)
		})
		const webmention = format.rels?.webmention?.[0]
		const pingback = format.rels?.pingback?.[0]
		if(!webmention && !pingback) {
			throw "no webmention and no pingback found?"
		}

		return {
			link: webmention ? webmention : pingback,
			type: webmention ? "webmention" : "pingback"
		}
	} catch(err) {
		console.warn(` -- whoops, failed to discover ${target}, why: ${err}`)
		return undefined
	}
}

module.exports = {
	discover
}
