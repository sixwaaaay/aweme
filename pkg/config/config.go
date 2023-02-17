// Package config provides the configuration for the application
package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config is the configuration for the application
type Config struct {
	// Release is the release mode of the application
	// if true, the application will run in release mode
	// zap logger will be used in production mode
	Release bool `yaml:"release"`
	JWT     struct {
		Secret    string   `yaml:"secret"`    // JwtSecret is the secret used to sign the JWT
		TTL       int64    `yaml:"ttl"`       // JwtTTL is the time to live for the JWT in seconds
		Whitelist []string `yaml:"whitelist"` // JwtWhitelist is the list of paths that can be accessed either with or without a JWT
	} `yaml:"jwt"` // JWT is the configuration for the JWT
	API struct {
		Prefix string `yaml:"prefix"` // ApiPrefix is the prefix for the API
	} `yaml:"api"` // API is the configuration for the API
	MySQL struct {
		Host     string `yaml:"host"`     // MySQLHost is the host of the MySQL database
		Port     int    `yaml:"port"`     // MySQLPort is the port of the MySQL database
		Username string `yaml:"username"` // MySQLUsername is the username of the MySQL database
		Password string `yaml:"password"` // MySQLPassword is the password of the MySQL database
		Database string `yaml:"database"` // MySQLDatabase is the database of the MySQL database
	} `yaml:"mysql"` // MySQL is the configuration for the MySQL database
}

// NewConfig returns a new Config
func NewConfig() (*Config, error) {
	path := "configs/config.yml" // default config file path
	if v := os.Getenv("CONFIG_FILE_PATH"); v != "" {
		path = v
	} // override config file path
	s, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c, err := parseConfig(s)
	return c, err
}

// parseConfig
// parses the binary config file data
// 1. expand environment variables in the config file
// 2. unmarshal the config file data to a Config struct
func parseConfig(s []byte) (*Config, error) {
	s = []byte(os.ExpandEnv(string(s))) // expand environment variables
	c := new(Config)
	err := yaml.Unmarshal(s, c) // unmarshal yaml to struct Config
	return c, err
}