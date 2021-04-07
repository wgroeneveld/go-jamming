
describe("webmention receiver validate tests", () => {

	const validhttpurl = "http://brainbaking.com/bla"
	const validhttpsurl = "https://brainbaking.com/blie"
	const urlfrominvaliddomain = "http://brainthe.bake/jup"
	const invalidurl = "lolzw"

	const { validate } = require('../../src/webmention/receive')


	test("is valid if source and target https urls", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded",
			body: {
				source: validhttpsurl + "1",
				target: validhttpsurl + "2"
			}
		})

		expect(result).toBe(true)
	})
	test("is NOT valid if target is a valid url but not form valid domain", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded",
			body: {
				source: validhttpsurl + "2",
				target: urlfrominvaliddomain
			}
		})

		expect(result).toBe(false)
	})
	test("is NOT valid if source and target are the same urls", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded",
			body: {
				source: validhttpsurl,
				target: validhttpsurl
			}
		})

		expect(result).toBe(false)
	})
	test("is valid if source and target http urls", () => {
		const result = validate({
			type: "application/x-www-form-urlencoded",
			body: {
				source: validhttpurl + "1",
				target: validhttpurl + "2"
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
