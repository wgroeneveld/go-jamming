# serve-my-jams 🥞

> A minimalistic jamstack-augmented microservice for webmentions etc

**Are you looking for a way to DO something with this?** See https://github.com/wgroeneveld/jam-my-stack !

This is a set of minimalistic [Koa-based](https://koajs.com/) microservices that aid you in your IndieWeb Jamstack coolness 😎 (name-dropping). While [jam-my-stack](https://github.com/wgroeneveld/jam-my-stack) is a set of scripts used to run at checkin-time, this is a dymamic service that handles requests. 

Inspect how it's used on https://brainbaking.com/ - usually, a `<link/>` in your `<head/>` suffices:

```
<link rel="webmention" href="https://jam.brainbaking.com/webmention" />
<link rel="pingback" href="https://webmention.io/webmention?forward=https://jam.brainbaking.com/webmention" />
```

If you want to support the older pingback protocol, you can leverage webmenton.io's forward capabilities. Although I developed this primarily because webmention.io is _not_ reliable - you've been warned. 

## What's in it?

### 1. Webmentions

#### 1.1 `POST /webmention`

Post a webmention. Includes a _lot_ of cross-checking and validating to guard against possible spam. See the [W3C WebMention spec](https://www.w3.org/TR/webmention/#sender-notifies-receiver) - or the source - for details.

Accepted form format: 

```
    POST /webmention-endpoint HTTP/1.1
    Host: aaronpk.example
    Content-Type: application/x-www-form-urlencoded

    source=https://waterpigs.example/post-by-barnaby&
    target=https://aaronpk.example/post-by-aaron
```

Will result in a `202 Accepted` - it handles things async. Stores in `.json` files in `data/domain`. 

#### 1.2 `GET /webmention/:domain/:token`

Retrieves a JSON array with relevant webmentions stored for that domain. The token should match. See `config.js` to fiddle with it yourself. Environment variables are supported, although I haven't used them yet. 

## TODOs

- `published` date is not well-formatted and blindly taken over from feed

