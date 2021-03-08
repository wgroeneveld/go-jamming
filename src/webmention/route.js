
const { validate, receive } = require('./receive')

function route(router) {
	router.post("webmention receive endpoint", "/webmention", async (ctx) => {
		if(!validate(ctx.request)) {
			ctx.throw(400, "malformed webmention request")
		}

		console.log(` OK: looks like a valid webmention: \n\tsource ${ctx.request.body.source}\n\ttarget ${ctx.request.body.target}`)
		// we do NOT await this on purpose.
		receive(ctx.request.body)

	    ctx.body = "Thanks, bro. Will process this webmention soon, pinky swear!";
	    ctx.status = 202
	});
}

module.exports = {
	route
}
