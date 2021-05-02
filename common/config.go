package common

import (
	"brainbaking.com/go-jamming/rest"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"strings"
	"time"
)

type Config struct {
	Port                     int      `json:"port"`
	Token                    string   `json:"token"`
	UtcOffset                int      `json:"utcOffset"`
	DataPath                 string   `json:"dataPath"`
	ConString                string   `json:"conString"`
	AllowedWebmentionSources []string `json:"allowedWebmentionSources"`
	Blacklist                []string `json:"blacklist"`
}

func (c *Config) IsBlacklisted(url string) bool {
	if !strings.HasPrefix(url, "http") {
		return false
	}
	domain := rest.Domain(url)
	return Includes(c.Blacklist, domain)
}

func (c *Config) Zone() *time.Location {
	return time.FixedZone("local", c.UtcOffset*60)
}

func (c *Config) missingKeys() []string {
	keys := []string{}
	if c.Port == 0 {
		keys = append(keys, "port")
	}
	if c.Token == "" {
		keys = append(keys, "token")
	}
	if c.ConString == "" {
		keys = append(keys, "conString")
	}
	if len(c.AllowedWebmentionSources) == 0 {
		keys = append(keys, "allowedWebmentionSources")
	}
	return keys
}

func (c *Config) IsAnAllowedDomain(domain string) bool {
	return Includes(c.AllowedWebmentionSources, domain)
}

func (c *Config) FetchDomain(url string) (string, error) {
	for _, domain := range c.AllowedWebmentionSources {
		if strings.Contains(url, domain) {
			return domain, nil
		}
	}
	return "", errors.New("no allowed domain found for url " + url)
}

func Configure() *Config {
	c := config()
	for _, domain := range c.AllowedWebmentionSources {
		log.Info().Str("allowedDomain", domain).Msg("Configured")
	}
	return c
}

func config() *Config {
	confData, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Warn().Msg("No config.json file found, reverting to defaults...")
		return defaultConfig()
	}

	conf := &Config{}
	err = json.Unmarshal(confData, conf)
	if err != nil {
		log.Warn().Msg("config.json malformed JSON, reverting to defaults...")
		return defaultConfig()
	}
	someMissingKeys := conf.missingKeys()
	if len(someMissingKeys) > 0 {
		log.Warn().Str("keys", strings.Join(someMissingKeys, ", ")).Msg("config.json is missing required keys, reverting to defaults...")
		return defaultConfig()
	}
	return conf
}

func defaultConfig() *Config {
	return &Config{
		Port:                     1337,
		Token:                    "miauwkes",
		UtcOffset:                60,
		ConString:                "mentions.db",
		AllowedWebmentionSources: []string{"brainbaking.com", "jefklakscodex.com"},
		Blacklist:                []string{"youtube.com"},
	}
}
