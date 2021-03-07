
const { receive, validate } = require('../../src/webmention/receive')

describe("validate tests", () => {

	const validhttpurl = "http://localhost/bla"
	const validhttpsurl = "https://localhost/blie"
	const invalidurl = "lolzw"

	test("is valid if source and target https urls", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded",
			body: {
				source: validhttpsurl,
				target: validhttpsurl
			}
		})

		expect(result).toBe(true)
	})
	test("is valid if source and target http urls", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded",
			body: {
				source: validhttpurl,
				target: validhttpurl
			}
		})

		expect(result).toBe(true)
	})
	test("is NOT valid if source is not a valid url", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded",
			body: {
				source: invalidurl,
				target: validhttpurl
			}
		})

		expect(result).toBe(false)
	})
	test("is NOT valid if target is not a valid url", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded",
			body: {
				source: validhttpurl,
				target: invalidurl
			}
		})

		expect(result).toBe(false)
	})
	test("is NOT valid if source is missing", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded",
			body: {
				target: validhttpurl
			}
		})

		expect(result).toBe(false)
	})
	test("is NOT valid if target is missing", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded",
			body: {
				source: validhttpurl
			}
		})

		expect(result).toBe(false)		
	})
	test("is NOT valid if no valid encoded form", () => {
		const result = validate({
			type: "ow-mai-got",
			body: {
				source: validhttpurl,
				target: validhttpurl
			}
		})

		expect(result).toBe(false)
	})
	test("is NOT valid if body is missing", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded"
		})

		expect(result).toBe(false)
	})

})

describe("receive webmention process tests", () => {

})