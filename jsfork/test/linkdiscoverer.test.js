
const { discover } = require('../src/linkdiscoverer')

describe("link discoverer", () => {

	test("discover 'unknown' if no link is present", async() => {
		const result = await discover("https://brainbaking.com/link-discover-test-none.html")
		expect(result).toEqual({
			type: "unknown"
		})
	})

	test("prefer webmentions over pingbacks if both links are present", async () => {
		const result = await discover("https://brainbaking.com/link-discover-bothtypes.html")
		expect(result).toEqual({
			link: "http://aaronpk.example/webmention-endpoint",
			type: "webmention"
		})
	})

	describe("discovers pingback links", () => {
		test("discover link if present in header", async () => {
			const result = await discover("https://brainbaking.com/pingback-discover-test.html")
			expect(result).toEqual({
				link: "http://aaronpk.example/pingback-endpoint",
				type: "pingback"
			})
		})

		test("discover link if sole entry somewhere in html", async () => {
			const result = await discover("https://brainbaking.com/pingback-discover-test-single.html")
			expect(result).toEqual({
				link: "http://aaronpk.example/pingback-endpoint-body",
				type: "pingback"
			})
		})

		test("use link in header if multiple present in html", async () => {
			const result = await discover("https://brainbaking.com/pingback-discover-test-multiple.html")
			expect(result).toEqual({
				link: "http://aaronpk.example/pingback-endpoint-header",
				type: "pingback"
			})
		})				
	})

	describe("discovers webmention links", () => {
		test("discover link if present in header", async () => {
			const result = await discover("https://brainbaking.com/link-discover-test.html")
			expect(result).toEqual({
				link: "http://aaronpk.example/webmention-endpoint",
				type: "webmention"
			})
		})

		test("discover link if sole entry somewhere in html", async () => {
			const result = await discover("https://brainbaking.com/link-discover-test-single.html")
			expect(result).toEqual({
				link: "http://aaronpk.example/webmention-endpoint-body",
				type: "webmention"
			})
		})

		test("use link in header if multiple present in html", async () => {
			const result = await discover("https://brainbaking.com/link-discover-test-multiple.html")
			expect(result).toEqual({
				link: "http://aaronpk.example/webmention-endpoint-header",
				type: "webmention"
			})
		})		
	})

})
