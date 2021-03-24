
const got = require('got')

const { send } = require('../../src/webmention/send')


describe("webmention send scenarios", () => {
	test("webmention send integration test that can send both webmentions and pingbacks", async () => {
		// jest.fn() gives unpredictable and unreadable output if unorderd calledWith... DIY!
		let posts = {}
		got.post = function(url, opts) {
			posts[url] = opts
		}

		// fetches index.xml
		await send("brainbaking.com", '2021-03-16T16:00:00.000Z')

		expect(Object.keys(posts).length).toBe(3)
		expect(posts["http://aaronpk.example/webmention-endpoint-header"]).toEqual({
			contentType: "x-www-form-urlencoded",
			form: {
				source: "https://brainbaking.com/notes/2021/03/16h17m07s14/",
				target: "https://brainbaking.com/link-discover-test-multiple.html"
			},
			retry: {
				limit: 5,
				methods: ["POST"]
			}
		})
		expect(posts["http://aaronpk.example/pingback-endpoint-body"]).toEqual({
			contentType: "text/xml",
			body: `<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://brainbaking.com/notes/2021/03/16h17m07s14/</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/pingback-discover-test-single.html</string></value>
		</param>
	</params>
</methodCall>`,
			retry: {
				limit: 5,
				methods: ["POST"]
			}
		})
		expect(posts["http://aaronpk.example/webmention-endpoint-body"]).toEqual({
			contentType: "x-www-form-urlencoded",
			form: {
				source: "https://brainbaking.com/notes/2021/03/16h17m07s14/",
				target: "https://brainbaking.com/link-discover-test-single.html"
			},
			retry: {
				limit: 5,
				methods: ["POST"]
			}
		})

	})
})