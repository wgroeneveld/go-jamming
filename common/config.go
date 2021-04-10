package common

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Port                        int
	Token                       string
	UtcOffset                   int
	DataPath                    string
	AllowedWebmentionSources    []string
	DisallowedWebmentionDomains []string
}

func (c *Config) ContainsDisallowedDomain(url string) bool {
	for _, domain := range c.DisallowedWebmentionDomains {
		if strings.Contains(url, domain) {
			return true
		}
	}
	return false
}

func (c *Config) IsAnAllowedDomain(url string) bool {
	for _, domain := range c.AllowedWebmentionSources {
		if domain == url {
			return true
		}
	}
	return false
}

func (c *Config) FetchDomain(url string) (string, error) {
	for _, domain := range c.AllowedWebmentionSources {
		if strings.Contains(url, domain) {
			return domain, nil
		}
	}
	return "", errors.New("no allowed domain found for url " + url)
}

func (c *Config) SetupDataDirs() {
	for _, domain := range c.AllowedWebmentionSources {
		os.MkdirAll(c.DataPath+"/"+domain, os.ModePerm)
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
		Port:                        port,
		Token:                       token,
		UtcOffset:                   60,
		DataPath:                    "data",
		AllowedWebmentionSources:    []string{"brainbaking.com", "jefklakscodex.com"},
		DisallowedWebmentionDomains: []string{"youtube.com"},
	}
	return
}
