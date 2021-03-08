
const config = require("../config")
const fsp = require('fs').promises

function validate(params) {
	return params.token === config.token &&
		config.allowedWebmentionSources.includes(params.domain)
}

async function load(domain) {
	const fileEntries = await fsp.readdir(`data/${domain}`, { withFileTypes: true });

	const files = await Promise.all(fileEntries.map(async (file) => {
		const contents = await fsp.readFile(`data/${domain}/${file.name}`, 'utf-8')
		return JSON.parse(contents)
	}));

	return files
}

module.exports = {
	validate,
	load
}
