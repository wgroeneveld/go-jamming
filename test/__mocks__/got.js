const fs = require('fs').promises

async function got(url) {
	const relativeUrl = url.replace('https://brainbaking.com/', '')
	const body = (await fs.readFile(`./test/__mocks__/${relativeUrl}`, 'utf8')).toString()

	let headers = {}
	try {
		headerFile = await fs.readFile(`./test/__mocks__/${relativeUrl.replace(".html", "")}-headers.json`, 'utf8')
		headers = JSON.parse(headerFile.toString())
	} catch {
	}
	
	return {
		headers,
		body
	}
}

async function gotPostMock(url, opts) {
}

got.post = gotPostMock

module.exports = got
