
const fs = require('fs');
const fsp = require('fs').promises;
const { rmdir } = require('./../utils')

jest.mock('got');
const md5 = require('md5')
const { receive } = require('../../src/webmention/receive')
const dumpdir = 'data'
const MockDate = require('mockdate')

describe("receive webmention process tests happy path", () => {

	beforeEach(() => {
		if(fs.existsSync(dumpdir)) {
			rmdir(dumpdir)
		}
		fs.mkdirSync(dumpdir)

		MockDate.set('2020-01-01')
	})

	test("receive saves a JSON file of indieweb-metadata if all is valid", async () => {
		await receive({
			source: "valid-indieweb-source.html",
			target: "valid-indieweb-target.html"
		})

		const filename = md5(`source=valid-indieweb-source.html,target=valid-indieweb-target.html`)
		const result = await fsp.readFile(`data/${filename}.json`, 'utf-8')
		const data = JSON.parse(result)

		expect(data).toEqual({
			author: {
				name: "Wouter Groeneveld",
				picture: "https://brainbaking.com//img/avatar.jpg"
			},
			content: "This is cool, I just found out about valid indieweb target - so cool...",
			source: "valid-indieweb-source.html",
			target: "valid-indieweb-target.html",
			published: "2021-03-06T12:41:00"
		})
	})

	test("receive saves a JSON file of non-indieweb-data such as title if all is valid", async () => {
		await receive({
			source: "valid-nonindieweb-source.html",
			target: "valid-indieweb-target.html"
		})

		const filename = md5(`source=valid-nonindieweb-source.html,target=valid-indieweb-target.html`)
		const result = await fsp.readFile(`data/${filename}.json`, 'utf-8')
		const data = JSON.parse(result)

		expect(data).toEqual({
			author: {
				name: "valid-nonindieweb-source.html",
			},
			content: "Diablo 2 Twenty Years Later: A Retrospective | Jefklaks Codex",
			source: "valid-nonindieweb-source.html",
			target: "valid-indieweb-target.html",
			published: "2020-01-01T01:00:00"
		})
	})

 	test("receive a target that does not point to the source does nothing", async () => {
		await receive({
			source: "valid-indieweb-source.html",
			target: "valid-indieweb-source.html"
		})

		const data = fs.readdirSync(dumpdir)
		expect(data.length).toBe(0)
	})


})
