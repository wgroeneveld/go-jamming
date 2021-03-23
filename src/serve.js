"use strict";

const Koa = require("koa");
const Logger = require("koa-logger");
const bodyParser = require('koa-body');
const koaRouter = require("koa-router");
const helmet = require("koa-helmet");
const { RateLimit } = require('koa2-ratelimit');

// koa docs: https://koajs.com/#application
const app = new Koa();
const router = new koaRouter();

// see https://www.npmjs.com/package/koa2-ratelimit, simple brute-force with helmet will suffice.
app.use(RateLimit.middleware({
  interval: { min: 15 },
  max: 100
}));
app.use(helmet());

// TODO not sure what to do on error yet
app.use(Logger());

// enable ctx.request.body parsing for x-www-form-urlencoded webmentions etc
app.use(bodyParser({
   multipart: true,
   urlencoded: true
}));

// route docs: https://github.com/koajs/router/blob/HEAD/API.md#module_koa-router--Router+get%7Cput%7Cpost%7Cpatch%7Cdelete%7Cdel
require("./webmention/route").route(router);
require("./pingback/route").route(router);
const config = require("./config");
config.setupDataDirs();

app.use(router.routes()).use(router.allowedMethods());

app.listen(config.port, config.host, () => {
	console.log(`Started localhost at port ${config.port}`)
});

