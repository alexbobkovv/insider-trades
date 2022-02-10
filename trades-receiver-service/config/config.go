package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type (
	Config struct {
		App      `yaml:"app"`
		Server   `yaml:"server"`
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
		URL string `env:"POSTGRES_URL"`
	}
)

const (
	yamlConfigPath = "."
	envConfigPath = "."
	envFileName = ".env"
)

func New() (*Config, error) {
	cfg := &Config{}

	//configFile, err := ioutil.ReadFile(configPath)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//if err := yaml.Unmarshal(configFile, cfg); err != nil {
	//	return nil, err
	//}

	viper.AddConfigPath(envConfigPath)
	viper.SetConfigFile(envFileName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(cfg.Postgres)

	if err != nil {
		return nil, err
	}

	log.Println("Ok")
	log.Println(cfg.Postgres.URL)

	os.Exit(-1)

	return cfg, nil
}