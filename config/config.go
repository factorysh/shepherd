package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Ttl map[string]time.Duration `yaml:"ttl"`
}

func New() *Config {
	return &Config{
		Ttl: make(map[string]time.Duration),
	}
}

func Read(path string) (*Config, error) {
	if path == "" {
		path = "/etc/shepherd.yml"
	}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
