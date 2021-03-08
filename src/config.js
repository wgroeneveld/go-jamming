
const { existsSync, mkdirSync } = require('fs')

const allowedWebmentionSources = [
	"brainbaking.com",
	"jefklakscodex.com"
]

function setupDataDirs() {
	allowedWebmentionSources.forEach(domain => {
		const dir = `data/${domain}`
		console.log(` -- configured for ${domain}`)
		if(!existsSync(dir)) {
			mkdirSync(dir, {
				recursive: true
			})
		}
	})
}


module.exports = {
	port: 4000,
	host: "localhost",

	utcOffset: 60,

	allowedWebmentionSources,
	setupDataDirs
}
