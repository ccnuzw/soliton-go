package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config is a wrapper around viper to provide configuration access.
type Config struct {
	v *viper.Viper
}

// NewConfig creates a new Config instance.
func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// Sensible defaults to keep generated projects runnable out-of-the-box.
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.dsn", "data.db")
	v.SetDefault("log.level", "info")
	
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// It's okay if config file is not found, we might rely on env vars
		fmt.Fprintln(os.Stderr, "Config file not found, relying on environment variables")
	}

	return &Config{v: v}, nil
}

// GetString returns a string value for the key.
func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

// GetInt returns an int value for the key.
func (c *Config) GetInt(key string) int {
	return c.v.GetInt(key)
}

// GetBool returns a bool value for the key.
func (c *Config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

// UnmarshalKey unmarshals a config section into a struct.
func (c *Config) UnmarshalKey(key string, rawVal interface{}) error {
	return c.v.UnmarshalKey(key, rawVal)
}
