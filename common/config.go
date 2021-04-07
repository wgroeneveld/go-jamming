
package common

import (
	"os"
	"strconv"
)

type Config struct {
	Port int
	Token string
	UtcOffset int
	AllowedWebmentionSources []string
	DisallowedWebmentionDomains []string
}

func Configure() (c *Config) {
	portstr := os.Getenv("PORT")
    port, err := strconv.Atoi(portstr)
    if err != nil {
		port = 1337
	}
	token := os.Getenv("TOKEN")
	if token == "" {
		token = "miauwkes"
	}

	c = &Config{
		Port: port,
		Token: token,
		UtcOffset: 60,
		AllowedWebmentionSources: []string{ "brainbaking.com", "jefklakscodex.com" },
		DisallowedWebmentionDomains: []string{ "youtube.com" },
	}
	return
}
