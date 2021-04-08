
package common

import (
	"os"
	"strconv"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Port int
	Token string
	UtcOffset int
	DataPath string
	AllowedWebmentionSources []string
	DisallowedWebmentionDomains []string
}

func (c *Config) FetchDomain(url string) (string, error) {
	for _, domain := range c.AllowedWebmentionSources {
		if strings.Index(url, domain) != -1 {
			return domain, nil
		}
	}
	return "", errors.New("no allowed domain found for url " + url)
}

func (c *Config) SetupDataDirs() {
	for _, domain := range c.AllowedWebmentionSources {
		os.MkdirAll(c.DataPath + "/" + domain, os.ModePerm)
		log.Info().Str("allowedDomain", domain).Msg("Configured")
	}
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
		DataPath: "data",
		AllowedWebmentionSources: []string{ "brainbaking.com", "jefklakscodex.com" },
		DisallowedWebmentionDomains: []string{ "youtube.com" },
	}
	return
}
