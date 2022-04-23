package common

import (
	"brainbaking.com/go-jamming/rest"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"io/fs"
	"io/ioutil"
	"strings"
	"time"
)

type Config struct {
	Port                     int      `json:"port"`
	Token                    string   `json:"token"`
	UtcOffset                int      `json:"utcOffset"`
	DataPath                 string   `json:"dataPath"`
	AllowedWebmentionSources []string `json:"allowedWebmentionSources"`
	Blacklist                []string `json:"blacklist"`
	Whitelist                []string `json:"whitelist"`
}

func (c *Config) IsBlacklisted(url string) bool {
	return isListedIn(url, c.Blacklist)
}

func (c *Config) IsWhitelisted(url string) bool {
	return isListedIn(url, c.Whitelist)
}

func isListedIn(url string, list []string) bool {
	if !strings.HasPrefix(url, "http") {
		return false
	}
	domain := rest.Domain(url)
	return Includes(list, domain)
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

// AddToBlacklist adds the given domain to the blacklist slice and persists to disk.
func (c *Config) AddToBlacklist(domain string) {
	c.Blacklist = addToList(domain, c.Blacklist)
	c.Save()
}

// AddToWhitelist adds the given domain to the whitelist slice and persists to disk.
func (c *Config) AddToWhitelist(domain string) {
	c.Whitelist = addToList(domain, c.Whitelist)
	c.Save()
}

func addToList(key string, arr []string) []string {
	for _, d := range arr {
		if d == key {
			return arr
		}
	}

	return append(arr, key)
}

func (c *Config) Save() {
	bytes, _ := json.Marshal(c) // we assume a correct internral state here
	err := ioutil.WriteFile("config.json", bytes, fs.ModePerm)
	if err != nil {
		log.Err(err).Msg("Unable to save config.json to disk!")
	}
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
		AllowedWebmentionSources: []string{"brainbaking.com", "jefklakscodex.com"},
		Blacklist:                []string{"youtube.com"},
		Whitelist:                []string{"brainbaking.com"},
	}
}
