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
	Postgres   PostgresConfig   `yaml:"postgres"`
	Redis      RedisConfig      `yaml:"redis"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
	Grafana    GrafanaConfig    `yaml:"grafana"`
	Login      LoginConfig      `yaml:"login"`
	Server     ServerConfig     `yaml:"server"`
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

type RedisConfig struct {
	Address string `yaml:"address"`
}

type PrometheusConfig struct {
	URL string `yaml:"url"`
}

type GrafanaConfig struct {
	URL      string `yaml:"url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type LoginConfig struct {
	Expiry int    `yaml:"expiry"` // in minutes
	Secret string `yaml:"secret"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}
