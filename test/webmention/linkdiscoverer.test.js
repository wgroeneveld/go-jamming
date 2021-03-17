
const { discover } = require('../../src/webmention/linkdiscoverer')

describe("link discoverer", () => {

	test("discover link if present in header", async () => {
		const result = await discover("https://brainbaking.com/link-discover-test.html")
		expect(result).toBe("http://aaronpk.example/webmention-endpoint")
	})

	test("discover nothing if no webmention link is present", async() => {
		const result = await discover("https://brainbaking.com/link-discover-test-none.html")
		expect(result).toBeUndefined()
	})

	test("discover link if sole entry somewhere in html", async () => {
		const result = await discover("https://brainbaking.com/link-discover-test-single.html")
		expect(result).toBe("http://aaronpk.example/webmention-endpoint-body")
	})

	test("use link in header if multiple present in html", async () => {
		const result = await discover("https://brainbaking.com/link-discover-test-multiple.html")
		expect(result).toBe("http://aaronpk.example/webmention-endpoint-header")
	})

})