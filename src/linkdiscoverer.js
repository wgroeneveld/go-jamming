const got = require('got')
const { mf2 } = require("microformats-parser");

const baseUrlOf = (url) => {
	if(url.match(/\//g).length <= 2) {
		return url
	}
	const split = url.split('/')
	return split[0] + '//' + split[2]
}

// see https://www.w3.org/TR/webmention/#sender-discovers-receiver-webmention-endpoint
async function discover(target) {
	try {
		const endpoint = await got(target)
		if(endpoint.headers.link?.indexOf("webmention") >= 0) {
			// e.g. Link: <http://aaronpk.example/webmention-endpoint>; rel="webmention"
			const link = endpoint.headers.link
				.split(";")[0]
				.replace("<" ,"")
				.replace(">", "")
			return {
				link,
				type: "webmention"
			}
		}

		const format = mf2(endpoint.body, {
			// this also complies with w3.org regulations: relative endpoint could be possible
			baseUrl: baseUrlOf(target)
		})
		const link = format.rels?.webmention?.[0]
		return {
			link,
			type: "webmention"
		}
	} catch(err) {
		console.warn(` -- whoops, failed to discover ${target}, why: ${err}`)
		return undefined
	}
}

module.exports = {
	discover
}
