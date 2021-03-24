describe("e2e tests", () => {

	jest.disableAutomock()
	jest.unmock('got')


	const { mf2 } = require("microformats-parser");
	const got = require("got");

	test.skip("microformat fiddling for non-indieweb sites", async () => {
		const html = (await got("https://kristienthoelen.be/2021/03/22/de-stadia-van-een-burn-out-in-welk-stadium-zit-jij/")).body
		const mf = mf2(html, {
			baseUrl: "https://kristienthoelen.be/"
		})

		//console.log(mf)

		const url = "https://kristienthoelen.be/wp-content/uploads/2021/03/burnoutbarometer.jpg"
		const occ = html.indexOf(url)
		const len = 100
		console.log(html.substring(occ - len, occ + url.length + len))

		<a[^>]+?' . $preg_target . '[^>]*>([^>]+?)</a>
	})

})
