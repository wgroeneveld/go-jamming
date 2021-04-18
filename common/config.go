package common

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"strings"
)

type Config struct {
	Port                        int      `json:"port"`
	Token                       string   `json:"token"`
	UtcOffset                   int      `json:"utcOffset"`
	DataPath                    string   `json:"dataPath"`
	ConString                   string   `json:"conString"`
	AllowedWebmentionSources    []string `json:"allowedWebmentionSources"`
	DisallowedWebmentionDomains []string `json:"disallowedWebmentionDomains"`
}

func (c *Config) missingKeys() []string {
	keys := []string{}
	if c.Port == 0 {
		keys = append(keys, "port")
	}
	if c.Token == "" {
		keys = append(keys, "token")
	}
	if c.DataPath == "" {
		keys = append(keys, "dataPath")
	}
	if c.ConString == "" {
		keys = append(keys, "conString")
	}
	if len(c.AllowedWebmentionSources) == 0 {
		keys = append(keys, "allowedWebmentionSources")
	}
	return keys
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

func Configure() *Config {
	conf := config()
	for _, domain := range conf.AllowedWebmentionSources {
		log.Info().Str("allowedDomain", domain).Msg("Configured")
	}
	return conf
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
		Port:                        1337,
		Token:                       "miauwkes",
		UtcOffset:                   60,
		DataPath:                    "data",
		ConString:                   "data/mentions.db",
		AllowedWebmentionSources:    []string{"brainbaking.com", "jefklakscodex.com"},
		DisallowedWebmentionDomains: []string{"youtube.com"},
	}
}
