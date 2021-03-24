
const pingbackReceiver = require('./receive')

function success(msg) {
	return `<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <params>
        <param>
            <value>
                <string>
                    ${msg}
                </string>
            </value>
        </param>
    </params>
</methodResponse>
`
}

function err(e) {
	console.err(` -- pingback receive went wrong: ${e.message}`)
	return `<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
    <fault>
        <value>
            <struct>
                <member>
                    <name>
                        faultCode
                    </name>
                    <value>
                        <int>
                            0
                        </int>
                    </value>
                </member>
                <member>
                    <name>
                        faultString
                    </name>
                    <value>
                        <string>
                        	${e.message}
                        </string>
                    </value>
                </member>
            </struct>
        </value>
    </fault>
</methodResponse>`
}

function route(router) {
	router.post("pingback receive endpoint", "/pingback", async (ctx) => {
		try {
			if(!pingbackReceiver.validate(ctx.request.body)) {
				throw "malformed pingback request"
			}

			// we do NOT await this on purpose.
			pingbackReceiver.receive(ctx.request.body)

		    ctx.status = 200
		    ctx.body = success("Thanks, bro. Will process this pingback soon, pinky swear!")
		} catch(e) {
			ctx.status = 200
			ctx.body = err(e)
		}
	});

    router.put("pingback send endpoint", "/pingback/:domain/:token", async (ctx) => {

    });
}

module.exports = {
	route
}
