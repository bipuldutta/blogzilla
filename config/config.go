package config

import (
	_ "embed"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

//go:embed config-local.yml
var configData []byte

type Config struct {
	Postgres    PostgresConfig    `yaml:"postgres"`
	DefaultUser DefaultUserConfig `yaml:"defaultuser"`
	Login       LoginConfig       `yaml:"login"`
	Server      ServerConfig      `yaml:"server"`
}

func NewConfig() *Config {
	// Parse the YAML data
	var conf Config
	if err := yaml.Unmarshal(configData, &conf); err != nil {
		log.Fatal(err)
	}
	fmt.Println(conf.Postgres.Host)
	return &conf
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// DefaultUserConfig a default admin user which will be created during database
// initialization stage when the service startsup for the first time.
// We recommend changing the password asap
type DefaultUserConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type LoginConfig struct {
	Expiry int    `yaml:"expiry"` // in minutes
	Secret string `yaml:"secret"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}
