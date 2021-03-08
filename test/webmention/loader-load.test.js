
const { load } = require('../../src/webmention/loader')
const fs = require('fs');
const fsp = require('fs').promises;
const { rmdir } = require('./../utils')
const dumpdir = 'data/brainbaking.com'

const exampleWebmention = {
	author: {
		name: "Wouter Groeneveld",
		picture: "https://brainbaking.com//img/avatar.jpg"
	},
	content: "This is cool, I just found out about valid indieweb target - so cool...",
	source: "https://coolness.com",
	target: "https://brainbaking.com/notes/2021/03/02h17m18s46/",
	published: "2021-03-06T12:41:00"
}

const exampleWebmention2 = {
	author: {
		name: "Jef Klakveld"
	},
	content: "Give it to me baby uhuh-uhuh white flies girls etc",
	source: "https://darkness.be",
	target: "https://brainbaking.com/about",
	published: "2021-03-06T12:41:00"
}

describe("webmention loading of existing json files tests", () => {
	beforeEach(() => {
		if(fs.existsSync(dumpdir)) {
			rmdir(dumpdir)
		}
		fs.mkdirSync(dumpdir)
	})

	test("return an array of webmentions from domain dir", async () => {
		await fsp.writeFile(`${dumpdir}/test.json`, JSON.stringify(exampleWebmention), 'utf-8')
		await fsp.writeFile(`${dumpdir}/test2.json`, JSON.stringify(exampleWebmention2), 'utf-8')

		const result = await load("brainbaking.com")
		expect(result.length).toBe(2)
	})

})