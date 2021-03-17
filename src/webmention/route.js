
const webmentionReceiver = require('./receive')
const webmentionLoader = require('./loader')
const webmentionSender = require('./send')

function route(router) {
	router.post("webmention receive endpoint", "/webmention", async (ctx) => {
		if(!webmentionReceiver.validate(ctx.request)) {
			ctx.throw(400, "malformed webmention request")
		}

		console.log(` OK: looks like a valid webmention: \n\tsource ${ctx.request.body.source}\n\ttarget ${ctx.request.body.target}`)
		// we do NOT await this on purpose.
		webmentionReceiver.receive(ctx.request.body)

	    ctx.body = "Thanks, bro. Will process this webmention soon, pinky swear!"
	    ctx.status = 202
	});

	router.put("webmention send endpoint", "/webmention/:domain/:token", async (ctx) => {
		if(!webmentionLoader.validate(ctx.params)) {
			ctx.throw(403, "access denied")
		}

		const since = ctx.request.query?.since
		console.log(` OK: someone wants to send mentions from domain ${ctx.params.domain} since ${since}`)
		// we do NOT await this on purpose.
		webmentionSender.send(ctx.params.domain, since)

		ctx.body = "Thanks, bro. Will send these webmentions soon, pinky swear!"
		ctx.status = 202
	})

	router.get("webmention get endpoint", "/webmention/:domain/:token", async (ctx) => {
		if(!webmentionLoader.validate(ctx.params)) {
			ctx.throw(403, "access denied")
		}

		console.log(` OK: someone wants a list of mentions at domain ${ctx.params.domain}`)
		const result = await webmentionLoader.load(ctx.params.domain)

		ctx.body = {
			status: 'success',
			json: result
		}
	})
}

module.exports = {
	route
}
