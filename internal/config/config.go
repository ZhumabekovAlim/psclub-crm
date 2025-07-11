package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	Auth struct {
		AccessSecret  string `yaml:"access_secret"`
		RefreshSecret string `yaml:"refresh_secret"`
		AccessTTL     int    `yaml:"access_ttl"`
		RefreshTTL    int    `yaml:"refresh_ttl"`
	} `yaml:"auth"`
}

func LoadConfig() *Config {
	f, err := os.Open("./config/config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
