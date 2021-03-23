
describe("pingback receive validation tests", () => {

	const { validate } = require('../../src/pingback/receive')

	test("not valid if malformed XML as body", () => {
		const result = validate("ola pola")
		expect(result).toBe(false)
	})

	test("not valid if methodName is not pingback.ping", () => {
		const result = validate(`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>ka.tsjing</methodName>
	<params>
		<param>
			<value><string>https://cool.site</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`)

		expect(result).toBe(false)
	})

	test("not valid if less than two parameters", () => {
		const result = validate(`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`)

		expect(result).toBe(false)
	})

	test("not valid if more than two parameters", () => {
		const xml = `<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://cool.site</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`

		expect(validate(xml)).toBe(false)
	})

	test("not valid if target is not in trusted domains from config", () => {
		const result = validate(`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://cool.site</string></value>
		</param>
		<param>
			<value><string>https://flashballz.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`)

		expect(result).toBe(false)
	})

	test("not valid if target is not http(s)", () => {
		const result = validate(`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://cool.site</string></value>
		</param>
		<param>
			<value><string>gemini://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`)

		expect(result).toBe(false)
	})

	test("not valid if source is not http(s)", () => {
		const result = validate(`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>gemini://cool.site</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`)

		expect(result).toBe(false)
	})

	test("is valid if pingback.ping and two http(s) parameters of which target is trusted", () => {
		const result = validate(`
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://cool.site</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/post/2021/03/cool-ness</string></value>
		</param>
	</params>
</methodCall>
			`)

		expect(result).toBe(true)
	})

})