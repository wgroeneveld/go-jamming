# go-jammin' 🥞

Go module `brainbaking.com/go-jamming`:

> A minimalistic Go-powered jamstack-augmented microservice for webmentions etc

✅️ **This is a fork of [serve-my-jams](https://github.com/wgroeneveld/serve-my-jams)**, the Node-powered original microservice, which is no longer being maintained. 

**Are you looking for a way to DO something with this?** See https://github.com/wgroeneveld/jam-my-stack !

This is a set of minimalistic Go-based microservices that aid you in your IndieWeb Jamstack coolness 😎 (name-dropping). While [jam-my-stack](https://github.com/wgroeneveld/jam-my-stack) is a set of scripts used to run at checkin-time, this is a dymamic service that handles requests. 

Inspect how it's used on https://brainbaking.com/ - usually, a `<link/>` in your `<head/>` suffices:

```
<link rel="webmention" href="https://jam.brainbaking.com/webmention" />
<link rel="pingback" href="https://jam.brainbaking.com/pingback" />
```

If you want to support the older pingback protocol, you can leverage webmenton.io's forward capabilities. Although I developed this primarily because webmention.io is _not_ reliable - you've been warned. 

## Building and running

Well, that's easy!

1. Build: `go build`
2. Run: `./go-jamming`
3. ???
4. Profit!

It's very much a fire-and-forget thing. Put it behind a reverse proxy such as nginx using something like this:

```
server {
        listen 443 ssl http2;
        listen [::]:443 ssl http2;

        server_name [your-domain];

        location / {
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $remote_addr;
                proxy_set_header Host $host;
                proxy_pass http://127.0.0.1:[your-port];
        }
    ssl_certificate /etc/letsencrypt/live/[your-domain]/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/[your-domain]/privkey.pem;
}
```

Create a very simple Linux system service that fires up the jam:

```
[Unit]
Description=Go-Jamming
After=network.target

[Service]
User=[myuser]
WorkingDirectory=/var/www/gojamming
ExecStart=/var/www/gojamming/go-jamming
SuccessExitStatus=0

[Install]
WantedBy=multi-user.target
```

Now install using `sudo systemctl enable/install gojamming` and you're done!

## Configuration

Place a `config.json` file in the same directory that looks like this: (below are the default values)

```json
{
  "port": 1337,
  "host": "localhost",
  "token": "miauwkes",
  "conString": "data/mentions.db",
  "utcOffset": 60,
  "allowedWebmentionSources":  [
    "brainbaking.com",
    "jefklakscodex.com"
  ],
  "disallowedWebmentionDomains":  [
    "youtube.com"
  ]
}
```

- port, host: http server params
- token, allowedWebmentionSources: see below, used for authentication
- disallowedWebmentionDomains: if an URL from that domain is encountered in your feed, ignore it. Does not send mentions to it. 
- utcOffset: offset in minutes for date processing, starting from UTC time.
- conString: file path to store all mentions and author avatars in a simple key/value store, based on [buntdb](https://github.com/tidwall/buntdb).

If a config file is missing, or required keys are missing, a warning will be generated and default values will be used instead. See `common/config.go`.

---

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

Will result in a `202 Accepted` - it handles things async. Stores in `.json` files in `[dataPath]/domain`. 

This also saves the author picture/avatar locally - if present in the microformat. It does _not_ resize images, however, if it's bigger than 5 MB, it falls back to a default one. 

Publication dates are sanitized and stored in `published`. They should be formatted in ISO8601. See [RFC3339](https://www.ietf.org/rfc/rfc3339.txt). 
If that is not the case, go-jamming falls back to the moment the mention was received. 

#### 1.2 `GET /webmention/:domain/:token`

Retrieves a JSON array with relevant webmentions stored for that domain. The token should match. See configuration to fiddle with it yourself. 

Example response:

```js
{
  "status": "success",
  "json": [
    {
      "author": {
        "name": "Jefklak",
        "picture": "/pictures/jefklakscodex.com"
      },
      "name": "Rainbow Six 3: Raven Shield - 17 Years Later",
      "content": "It’s amazing that the second disk is still readable by my Retro WinXP machine. It has been heavily abused in 2003 and the years after that. Rainbow Six' third installment, Raven Shield (or simply RvS), is quite a departure from the crude looking Rogu...",
      "published": "2020-11-01",
      "url": "https://jefklakscodex.com/articles/retrospectives/raven-shield-17-years-later/",
      "type": "mention",
      "source": "https://jefklakscodex.com/articles/retrospectives/raven-shield-17-years-later/",
      "target": "https://brainbaking.com/post/2020/10/building-a-core2duo-winxp-retro-pc/"
    }
  ]
}
```

A few remarks:

- `picture`: Author picture paths are relative to the jamming server since they're locally stored. 
- `published`: This is not processed and simply taken over from the microformat.
- `target` is your domain, `source` is... well... the source. 
- `content`: Does not contain HTML. Automatically capped at 250 characters if needed.

#### 1.3 `PUT /webmention/:domain/:token`

Sends out **both webmentions and pingbacks**, based on the domain's `index.xml` RSS feed, and optionally, a `since` request query parameter that is supposed to be a string, fed through [Dayjs](https://day.js.org/) to format. (e.g. `2021-03-16T16:00:00.000Z`). 

This does a couple of things:

1. Fetch RSS entries (since last sent link x, or everything)
2. Find outbound `href`s (starting with `http`)
3. Check if those domains have a `webmention` link endpoint installed, according to the w3.org rules. If not, check for a `pingback` endpoint. If not, bail out.
4. If webmention/pingback found: `POST` for each found href with `source` the own domain and `target` the outbound link found in the RSS feed, using either XML or form data according to the protocol. 

As with the `POST` call, will result in a `202 Accepted` and handles things async/in parallel. 

**Does this thing take updates into account**?

Yes and no. It checks the `<link/>` tag to see if there's a new post since mentions were last sent. If a new link is discovered, it will send out those. 

This means if you made changes in-between, and they appear in the RSS feed as recent items, it will get resend. 

**Do I have to provide a ?source= parameter each time**?

No. The server will automatically store the latest push, and if it's called again, it will not send out anything if nothing more recent was found in your RSS feed based on the last published link. Providing the parameter merely lets you override the behavior.

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

## Troubleshooting

Run in verbose mode: use `-verbose`. This also logs debug info. Structured JSON is outputted through os.Stderr - which is usually `/var/log/syslog`. 

If rolling files in a separate location is required, [lumberjack](https://github.com/natefinch/lumberjack) could be added in `main.go`.

There's a **rate limiting** system implemented with a rate limit of 5 requests per second and a maximum burst rate of 10. 
That's pretty flexible. I have not taken the trouble to put this into the config, it should do in most cases. If you get a `429 too many requests`, you've hit the limiter. 
A separate goroutine cleans up ips each 2 minutes, the TTL is 5 minutes. See `limiter.go`. 

Database migrations are run using the `-migrate` flag. 

## TODOs

- Pictures are bound to domain names only. That means `brid.gy` will net a single picture. Perhaps the combination domain + user would be more appropriate?

