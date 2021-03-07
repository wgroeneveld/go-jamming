
function isValidUrl(url) {
	return url !== undefined &&
		(url.startsWith("http://") || url.startsWith("https://"))
}

/**
 https://www.w3.org/TR/webmention/#sender-notifies-receiver
 example:
		POST /webmention-endpoint HTTP/1.1
		Host: aaronpk.example
		Content-Type: application/x-www-form-urlencoded

		source=https://waterpigs.example/post-by-barnaby&
		target=https://aaronpk.example/post-by-aaron


		HTTP/1.1 202 Accepted
*/	 
function validate(request) {
	return request.type === "application/x-www-form-urlencoded" &&
		request.body !== undefined &&
		isValidUrl(request?.body.source) &&
		isValidUrl(request?.body.target)
}

async function receive(body) {
	// do stuff with it
} 

module.exports = {
	receive,
	validate
}
