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

Receive a webmention. Includes a _lot_ of cross-checking and validating to guard against possible spam. See the [W3C WebMention spec](https://www.w3.org/TR/webmention/#sender-notifies-receiver) - or the source - for details.

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

#### 1.3 `PUT /webmention/:domain/:token`

Sends out **both webmentions and pingbacks**, based on the domain's `index.xml` RSS feed, and optionally, a `since` request query parameter that is supposed to be a string, fed through [Dayjs](https://day.js.org/) to format. (e.g. `2021-03-16T16:00:00.000Z`). 

This does a couple of things:

1. Fetch RSS entries (since x, or everything)
2. Find outbound `href`s (starting with `http`)
3. Check if those domains have a `webmention` link endpoint installed, according to the w3.org rules. If not, check for a `pingback` endpoint. If not, bail out.
4. If webmention/pingback found: `POST` for each found href with `source` the own domain and `target` the outbound link found in the RSS feed, using either XML or form data according to the protocol. 

As with the `POST` call, will result in a `202 Accepted` and handles things async/in parallel. 

**Does this thing take updates into account**?

Yes and no. It checks the `<pubDate/>` `<item/>` RSS tag by default, but if a `<time datetime="..."/>` tag is present in the `<description/>`, it treats that date as the "last modified" date. There is no such thing in the RSS 2.0 W3.org specs, so I had to come up with my own hacks! Remember that if you want this to work, you also need to include a time tag in your RSS feed (e.g. `.Lastmod` gitinfo in Hugo). 

### 2. Pingbacks

Pingbacks are in here for two reasons:

1. I wanted to see how difficult it was to implement them. Turns out to be almost exactly the same as webmentions. This means the "new" W3 standards for webmentions are just as crappy as pingbacks... What's the difference between a form POST and an XML POST? Form factor?
2. Much more blogs (Wordpress-alike) support only pingbacks. 

#### 2.1 `POST /pingback`

Receive a pingback. Includes a _lot_ of cross-checking and validating to guard against possible spam. Internally, converts it into a webmention and processes it just like that.

Accepted XML body: 

```
<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
    <methodName>pingback.ping</methodName>
    <params>
        <param>
            <value><string>https://brainbaking.com/kristien.html</string></value>
        </param>
        <param>
            <value><string>https://kristienthoelen.be/2021/03/22/de-stadia-van-een-burn-out-in-welk-stadium-zit-jij/</string></value>
        </param>
    </params>
</methodCall>
```

Will result in a `200 OK` - that returns XML according to [The W3 pingback XML-RPC spec](https://www.hixie.ch/specs/pingback/pingback#refsXMLRPC). Processes async. 

#### 2.2 Sending pingbacks

Happens automatically through `PUT /webmention/:domain/:token`! Links that are discovered as `rel="pingback"` that **do not** already have a webmention link will be processed as XML-RPC requests to be send. 


## TODOs

- `published` date is not well-formatted and blindly taken over from feed
- Implement a Brid.gy-like system that converts links from domains in the config found on [public Mastodon timelines](https://docs.joinmastodon.org/methods/timelines/) into webmentions. (And check if it's ok to only use the public line)
