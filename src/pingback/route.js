
const pingbackReceiver = require('./receive')
const log = require('pino')()

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
	log.error(e, 'pingback receive went wrong')
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
                        	Sorry pal. Malformed request? Or something else, who knows...
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
}

module.exports = {
	route
}
