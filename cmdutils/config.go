package cmdutils

import (
	"os"
)

var configDefaults = map[string]string{
	"APP_PORT": "7447",
}

// Config is the store for config data
type Config struct {
	data map[string]string
}

// NewConfig will get all environment variables and set defaults
func NewConfig() Config {

	c := Config{map[string]string{}}

	for k, d := range configDefaults {
		v := os.Getenv(k)
		c.set(k, v, d)
	}

	return c

}

func (c *Config) set(key, val, fallback string) {

	if len(val) == 0 {
		c.data[key] = fallback
	} else {
		c.data[key] = val
	}

}

// Get will retrieve a config setting
func (c Config) Get(key string) string {

	return c.data[key]

}
