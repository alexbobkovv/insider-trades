package config

import (
	"gopkg.in/yaml.v3"
	_ "gopkg.in/yaml.v3"
	"io/ioutil"
)

type (
	Config struct {
		App `yaml:"app"`
		Server `yaml:"server"`
		Postgres `yaml:"postgres"`
	}

	App struct {
		Name string `yaml:"name"`
		Version string `yaml:"version"`
	}

	Server struct {
		Port string `yaml:"port"`
	}

	Postgres struct {
		PoolMax int `yaml:"pool-max"`
		URL string `yaml:"url"`
	}
)

const configPath = "configs/config.yml"

func New() (*Config, error) {
	cfg := &Config{}

	configFile, err := ioutil.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(configFile, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}