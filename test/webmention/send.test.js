
const got = require('got')

const { send } = require('../../src/webmention/send')


describe("webmention send scenarios", () => {
	test("webmention send integration test", async () => {
		got.post = jest.fn()

		// fetches index.xml
		await send("brainbaking.com", '2021-03-16T16:00:00.000Z')

		expect(got.post).toHaveBeenCalledTimes(2)
		expect(got.post).toHaveBeenCalledWith("http://aaronpk.example/webmention-endpoint-header", {
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
		expect(got.post).toHaveBeenCalledWith("http://aaronpk.example/webmention-endpoint-body", {
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