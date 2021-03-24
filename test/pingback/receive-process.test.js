
const fs = require('fs');
const fsp = require('fs').promises;
const { rmdir } = require('./../utils')

jest.mock('got');
const md5 = require('md5')
const { receive } = require('../../src/pingback/receive')
const dumpdir = 'data/brainbaking.com'

describe("receive pingback process tests happy path", () => {

	beforeEach(() => {
		if(fs.existsSync(dumpdir)) {
			rmdir(dumpdir)
		}
		fs.mkdirSync(dumpdir, {
			recursive: true
		})
	})

	function asFilename(body) {
		return `${dumpdir}/` + md5(`source=${body.source},target=${body.target}`)
	}

	test("receive a pingback processes it just like a webmention", async () => {
		const body = {
			source: "https://brainbaking.com/valid-bridgy-twitter-source.html",
			target: "https://brainbaking.com/post/2021/03/the-indieweb-mixed-bag"
		}

		await receive(`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>${body.source}</string></value>
		</param>
		<param>
			<value><string>${body.target}</string></value>
		</param>
	</params>
</methodCall>
			`)

		const result = await fsp.readFile(`${asFilename(body)}.json`, 'utf-8')
		const data = JSON.parse(result)
		expect(data.content).toContain("Recommended read:")
	})

})
