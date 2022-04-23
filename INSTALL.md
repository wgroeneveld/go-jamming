# Go-Jamming installation guide

Back to the main [README.md](https://github.com/wgroeneveld/go-jamming/blob/master/README.md)

## 1. Installing

Download the latest binary version from the [GitHub releases page](https://github.com/wgroeneveld/go-jamming/releases). This is a **single binary** and installing it is just a matter of copy-pasting it to your sever! Simply execute with `./go-jamming`. 

### Compiling it yourself (optional)

If your target OS is not listed, you can build it yourself with one simple command: `go build`. Go 1.16+ is required, see `go.mod` file.

## 2. Configuring

Place a `config.json` file in the same directory that looks like this: (below are the default values)

```json
{
  "port": 1337,
  "host": "localhost",
  "token": "miauwkes",
  "utcOffset": 60,
  "allowedWebmentionSources":  [
    "brainbaking.com",
    "jefklakscodex.com"
  ],
  "blacklist":  [
    "youtube.com"
  ]
}
```

- port, host: http server params
- token, allowedWebmentionSources: see below, used for authentication
- blacklist: blacklist domains from which we do NOT send to or accept mentions from. 
- utcOffset: offset in minutes for date processing, starting from UTC time.

If a config file is missing, or required keys are missing, a warning will be generated and default values will be used instead. See `common/config.go`.

To keep things simple, the file path to store all mentions and author avatars in a simple key/value store is hardcoded and set to:

- mentions.db (in working dir) for approved mentions
- mentions_toapprove.db (in working dir) for mentions in moderation.

The database is based on [buntdb](https://github.com/tidwall/buntdb). If the files do not exist, they will simply be created.


## 3. Reverse proxy

Put it behind a reverse proxy such as nginx using something like this:

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

## 4. Linux systemd service

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

## 5. Configuring your templates

Read more on how to integrate this in for example Hugo on https://brainbaking.com/post/2021/05/beyond-webmention-io/
