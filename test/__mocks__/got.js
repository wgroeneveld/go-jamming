const fs = require('fs').promises

async function got(url) {
	const relativeUrl = url.replace('https://brainbaking.com/', '')
	const body = await fs.readFile(`./test/__mocks__/${relativeUrl}`, 'utf8')
	return {
		body
	}
}

module.exports = got
