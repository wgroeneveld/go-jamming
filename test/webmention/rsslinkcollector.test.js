
const { collect } = require('../../src/webmention/rsslinkcollector')
const fs = require('fs').promises
const dayjs = require('dayjs')

describe("collect RSS links of articles since certain period", () => {

	let xml = ''
	beforeEach(async () => {
		xml = (await fs.readFile('./test/__mocks__/samplerss.xml')).toString()
	})

	test("collect should not contain hrefs from blocked domains", () => {
		const collected = collect(xml, dayjs('2021-03-10T00:00:00.000Z').toDate())

		// test case: 
		// contains youtube.com/cool link
		const last = collected[collected.length - 1]
		expect(last.hrefs).toEqual([
			"https://dog.estate/@eli_oat",
			"https://twitter.com/olesovhcom/status/1369478732247932929",
			"/about"
		])

	})

	test("collect should not contain hrefs that point to images", () => {
		const collected = collect(xml, dayjs('2021-03-14T00:00:00.000Z').toDate())

		// test case: 
		// contains e.g. https://chat.brainbaking.com/media/6f8b72ca-9bfb-460b-9609-c4298a8cab2b/EuropeBattle%202021-03-14%2016-20-36-87.jpg
		const last = collected[collected.length - 1]
		expect(last.hrefs).toEqual([
			"/about"
		])
	})

	test("collects if time tag found in content that acts as an update stamp", async () => {
		// sample item: pubDate 2021-03-16, timestamp updated: 2021-03-20
		xml = (await fs.readFile('./test/__mocks__/samplerss-updated-timestamp.xml')).toString()

		const collected = collect(xml, dayjs('2021-03-19').toDate())
		expect(collected.length).toBe(1)
	})

	test("does not collect if time tag found in content but still older than since", async () => {
		// sample item: pubDate 2021-03-16, timestamp updated: 2021-03-20
		xml = (await fs.readFile('./test/__mocks__/samplerss-updated-timestamp.xml')).toString()

		const collected = collect(xml, dayjs('2021-03-21').toDate())
		expect(collected.length).toBe(0)
	})

	test("collects nothing if date in future and since nothing new in feed", () => {
		const collected = collect(xml, dayjs().add(7, 'day').toDate())
		expect(collected.length).toEqual(0)
	})

	test("collect latest x links when a since parameter is provided", () => {
		const collected = collect(xml, dayjs('2021-03-15T00:00:00.000Z').toDate())
		expect(collected.length).toEqual(3)

		const last = collected[collected.length - 1]
		expect(last.link).toBe("https://brainbaking.com/notes/2021/03/15h14m43s49/")
		expect(last.hrefs).toEqual([
			"http://replit.com",
			"http://codepen.io",
			"https://kuleuven-diepenbeek.github.io/osc-course/ch1-c/intro/",
			"/about"
		])
	})

	test("collect every external link without a valid since date", () => {
		const collected = collect(xml)
		expect(collected.length).toEqual(141)

		const first = collected[0]
		expect(first.link).toBe("https://brainbaking.com/notes/2021/03/16h17m07s14/")
		expect(first.hrefs).toEqual([
			"https://fosstodon.org/@celia",
			"https://fosstodon.org/@kev",
			"/about"
		])
	})

})