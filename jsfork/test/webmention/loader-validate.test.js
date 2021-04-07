
describe("webmention loader validate tests", () => {

	const { validate } = require('../../src/webmention/loader')
	const config = require('../../src/config')

	test("is invalid if token not the same", () => {
		const result = validate({
			token: "drie roze olifanten hopla in de lucht",
			domain: config.allowedWebmentionSources[0]
		})

		expect(result).toBe(false)
	})

	test("is invalid if domain not the list of known domains", () => {
		const result = validate({
			token: config.token,
			domain: "woozaas.be"
		})

		expect(result).toBe(false)
	})

	test("is valid if domain and token matching", () => {
		const result = validate({
			token: config.token,
			domain: config.allowedWebmentionSources[0]
		})

		expect(result).toBe(true)
	})
})
