
const got = require('got')

async function sendPingbackToEndpoint(endpoint, source, target) {
	const body = `<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>${source}</string></value>
		</param>
		<param>
			<value><string>${target}</string></value>
		</param>
	</params>
</methodCall>`
	await got.post(endpoint, {
		contentType: "text/xml",
		body,
		retry: {
			limit: 5,
			methods: ["POST"]
		}
	})
	console.log(` OK: pingback@${endpoint}, sent: source ${source}, target ${target}`)
}

module.exports = {
	sendPingbackToEndpoint
}
