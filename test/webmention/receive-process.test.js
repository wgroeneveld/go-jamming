
const fs = require('fs');
const fsp = require('fs').promises;
const { rmdir } = require('./../utils')

jest.mock('got');
const md5 = require('md5')
const { receive } = require('../../src/webmention/receive')
const dumpdir = 'data/brainbaking.com'
const MockDate = require('mockdate')

describe("receive webmention process tests happy path", () => {

	beforeEach(() => {
		if(fs.existsSync(dumpdir)) {
			rmdir(dumpdir)
		}
		fs.mkdirSync(dumpdir)

		MockDate.set('2020-01-01')
	})

	function asFilename(body) {
		return `${dumpdir}/` + md5(`source=${body.source},target=${body.target}`)
	}

	test("receive saves a JSON file of indieweb-metadata if all is valid", async () => {
		const body = {
			source: "https://brainbaking.com/valid-indieweb-source.html",
			target: "https://brainbaking.com/valid-indieweb-target.html"
		}
		await receive(body)

		const result = await fsp.readFile(`${asFilename(body)}.json`, 'utf-8')
		const data = JSON.parse(result)

		expect(data).toEqual({
			author: {
				name: "Wouter Groeneveld",
				picture: "https://brainbaking.com//img/avatar.jpg"
			},
			content: "This is cool, I just found out about valid indieweb target - so cool...",
			source: body.source,
			target: body.target,
			published: "2021-03-06T12:41:00"
		})
	})

	test("receive saves a JSON file of indieweb-metadata with summary as content if present", async () => {
		const body = {
			source: "https://brainbaking.com/valid-indieweb-source-with-summary.html",
			target: "https://brainbaking.com/valid-indieweb-target.html"
		}
		await receive(body)

		const result = await fsp.readFile(`${asFilename(body)}.json`, 'utf-8')
		const data = JSON.parse(result)

		expect(data).toEqual({
			author: {
				name: "Wouter Groeneveld",
				picture: "https://brainbaking.com//img/avatar.jpg"
			},
			content: "This is cool, this is a summary!",
			source: body.source,
			target: body.target,
			published: "2021-03-06T12:41:00"
		})
	})

	test("receive saves a JSON file of non-indieweb-data such as title if all is valid", async () => {
		const body = {
			source: "https://brainbaking.com/valid-nonindieweb-source.html",
			target: "https://brainbaking.com/valid-indieweb-target.html"
		}
		await receive(body)

		const result = await fsp.readFile(`${asFilename(body)}.json`, 'utf-8')
		const data = JSON.parse(result)

		expect(data).toEqual({
			author: {
				name: "https://brainbaking.com/valid-nonindieweb-source.html",
			},
			content: "Diablo 2 Twenty Years Later: A Retrospective | Jefklaks Codex",
			source: body.source,
			target: body.target,
			published: "2020-01-01T01:00:00"
		})
	})

 	test("receive a target that does not point to the source does nothing", async () => {
 		const body = {
			source: "https://brainbaking.com/valid-indieweb-source.html",
			target: "https://brainbaking.com/valid-indieweb-source.html"
		}
		await receive(body)

		const data = fs.readdirSync(dumpdir)
		expect(data.length).toBe(0)
	})

 	test("receive a source that does not exist should also delete older webmention files", async () => {
 		const body = {
			source: "https://wubanga2001.boom/lolz",
			target: "https://brainbaking.com/valid-indieweb-source.html"
		}

		await fsp.writeFile(`${asFilename(body)}.json`, JSON.stringify({ lolz: "aha" }), 'utf-8')
		await receive(body)

		const data = fs.readdirSync(dumpdir)
		expect(data.length).toBe(0)
	})


})
