package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		App      `yaml:"app"`
		Server   `yaml:"server"`
		Postgres `yaml:"postgres"`
		Logger   `yaml:"zap"`
		RabbitMQ `yaml:"rabbit_mq"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	Server struct {
		Port string `yaml:"port"`
	}

	Postgres struct {
		URL string `mapstructure:"POSTGRES_URL"`
	}

	Logger struct {
		Level    string `yaml:"level"`
		Format   string `yaml:"format"`
		Filepath string `yaml:"filepath"`
	}

	RabbitMQ struct {
		AmqpURL string `mapstructure:"AMQP_URL"`
	}
)

const (
	yamlFileName       = "config.yml"
	yamlConfigPath     = "./config/"
	localEnvFileName   = ".env"
	localEnvConfigPath = "."
)

func New() (*Config, error) {
	cfg := &Config{}

	err := cfg.parseYAML(yamlFileName, yamlConfigPath)

	if err != nil {
		return nil, err
	}

	err = cfg.parseEnv(localEnvFileName, localEnvConfigPath)

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

	err = viper.Unmarshal(&c.Postgres)
	if err != nil {
		return err
	}

	viper.AutomaticEnv()
	err = viper.Unmarshal(&c.Postgres)
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&c.RabbitMQ)
	if err != nil {
		return err
	}

	return nil
}
