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
}

func LoadConfig() *Config {
	f, err := os.Open("C:\\Users\\alimz\\GolandProjects\\clean_mobile_app\\config\\config.yaml")
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
