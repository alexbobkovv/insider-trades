package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		App    `yaml:"app"`
		Server `yaml:"server"`
		Telegram
		Logger `yaml:"logger"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	Telegram struct {
		BotToken string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	}

	Server struct {
		Port string `yaml:"port"`
	}

	Logger struct {
		Level    string `yaml:"level"`
		Format   string `yaml:"format"`
		Filepath string `yaml:"filepath"`
	}
)

const (
	yamlFileName   = "config.yml"
	yamlConfigPath = "./config/"
	envFileName    = ".env"
	envConfigPath  = "."
)

func New() (*Config, error) {
	cfg := &Config{}

	err := cfg.parseYAML(yamlFileName, yamlConfigPath)

	if err != nil {
		return nil, err
	}

	err = cfg.parseEnv(envFileName, envConfigPath)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) parseYAML(fileName, filePath string) error {

	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filePath)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	err = viper.Unmarshal(&c)

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) parseEnv(fileName, filePath string) error {

	viper.SetConfigName(fileName)
	viper.SetConfigType("env")
	viper.AddConfigPath(filePath)

	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	// err = viper.Unmarshal(&c.Telegram)
	//
	// if err != nil {
	// 	return err
	// }
	//
	// viper.AutomaticEnv()
	err = viper.Unmarshal(&c.Telegram)

	if err != nil {
		return err
	}

	return nil
}
