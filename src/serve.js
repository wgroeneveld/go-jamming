"use strict";

const Koa = require("koa");
const Logger = require("koa-logger");
const bodyParser = require('koa-body');
const koaRouter = require("koa-router");
const helmet = require("koa-helmet");

// koa docs: https://koajs.com/#application
const app = new Koa();
const router = new koaRouter();

// TODO not sure what to do on error yet
app.use(Logger());

// enable ctx.request.body parsing for x-www-form-urlencoded webmentions etc
app.use(bodyParser({
   multipart: true,
   urlencoded: true
}));


// route docs: https://github.com/koajs/router/blob/HEAD/API.md#module_koa-router--Router+get%7Cput%7Cpost%7Cpatch%7Cdelete%7Cdel
require("./webmention/route").route(router)

app.use(helmet());
app.use(router.routes()).use(router.allowedMethods());

const port = process.env.PORT || 4000

app.listen(port, "localhost", () => {
	console.log(`Started localhost at port ${port}`)
});

