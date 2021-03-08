
const got = require('got')
const config = require('./../config')

function isValidUrl(url) {
	return url !== undefined &&
		(url.startsWith("http://") || url.startsWith("https://"))
}

function isValidDomain(url) {
	return config.allowedWebmentionSources.some(domain => {
		return url.indexOf(domain) !== -1
	})
}

/**
Remember, TARGET is own domain, SOURCE is the article to process
 https://www.w3.org/TR/webmention/#sender-notifies-receiver
 example:
		POST /webmention-endpoint HTTP/1.1
		Host: aaronpk.example
		Content-Type: application/x-www-form-urlencoded

		source=https://waterpigs.example/post-by-barnaby&
		target=https://aaronpk.example/post-by-aaron


		HTTP/1.1 202 Accepted
*/	 
function validate(request) {
	return request.type === "application/x-www-form-urlencoded" &&
		request.body !== undefined &&
		isValidUrl(request?.body?.source) &&
		isValidUrl(request?.body?.target) &&
		request?.body?.source !== request?.body?.target &&
		isValidDomain(request?.body?.source)
}

function processSourceBody(body, target) {
	if(body.indexOf(target) === -1) {
		return
	}
}

async function receive(body) {
	try {
		await got(body.target)
	} catch(unknownTarget) {
		return
	}

	try {
		const src = await got(body.source)
		processSourceBody(src.body, body.target)
	} catch(unknownSource) {
		return
	}
} 

module.exports = {
	receive,
	validate
}
