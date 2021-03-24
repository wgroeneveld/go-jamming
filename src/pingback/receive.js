
const webmentionReceiver = require('./../webmention/receive')
const config = require('./../config')
const parser = require("fast-xml-parser")

/**
See https://www.hixie.ch/specs/pingback/pingback#refsXMLRPC
---
Sample XML:
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://brainbaking.com/kristien.html</string></value>
		</param>
		<param>
			<value><string>https://kristienthoelen.be/2021/03/22/de-stadia-van-een-burn-out-in-welk-stadium-zit-jij/</string></value>
		</param>
	</params>
</methodCall>
*/
const isValidDomain = (url) => {
	return config.allowedWebmentionSources.some(domain => {
		return url.indexOf(domain) !== -1
	})
}

function validate(body) {
	const xml = parser.parse(body)
	if(!xml) return false
	if(!xml.methodCall || xml.methodCall.methodName !== "pingback.ping") return false
	if(!xml.methodCall.params || !xml.methodCall.params.param || xml.methodCall.params.param.length !== 2) return false
	if(!xml.methodCall.params.param.every(param => param?.value?.string?.startsWith('http'))) return false
	if(!isValidDomain(xml.methodCall.params.param[1].value.string)) return false
	return true
}

// we treat a pingback as a webmention. 
// Wordpress pingback processing source: https://developer.wordpress.org/reference/classes/wp_xmlrpc_server/pingback_ping/
async function receive(body) {
	const xml = parser.parse(body)
	const webmentionBody = {
		source: xml.methodCall.params.param[0].value.string,
		target: xml.methodCall.params.param[1].value.string
	}

    console.log(` OK: looks like a valid pingback: \n\tsource ${webmentionBody.source}\n\ttarget ${webmentionBody.target}`)
	await webmentionReceiver.receive(webmentionBody)
}

module.exports = {
	receive,
	validate
}
