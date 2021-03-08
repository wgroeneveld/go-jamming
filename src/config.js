
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
	port: process.env.PORT || 1337,
	host: "localhost",
	token: process.env.TOKEN || "miauwkes",

	utcOffset: 60,

	allowedWebmentionSources,
	setupDataDirs
}
