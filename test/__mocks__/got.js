const fs = require('fs').promises

async function got(url) {
	const body = await fs.readFile(`./test/__mocks__/${url}`, 'utf8')
	return {
		body
	}
}

module.exports = got
