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
		Port         string `yaml:"port"`
		ReceiverPath string `mapstructure:"RECEIVER_PATH"`
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
		AmqpURL      string `mapstructure:"AMQP_URL"`
		Exchange     string `yaml:"exchange"`
		Durable      bool   `yaml:"durable"`
		QueueName    string `yaml:"queueName"`
		RoutingKey   string `yaml:"routingKey"`
		ConsumerName string `yaml:"consumerName"`
	}
)

const (
	defaultYamlConfigName     = "config.yml"
	defaultYamlConfigPath     = "./config/"
	defaultLocalEnvConfigName = ".env"
	defaultLocalEnvConfigPath = "."
)

func New(configPath, configName string) (*Config, error) {
	if configPath == "" {
		configPath = defaultYamlConfigPath
	}
	if configName == "" {
		configName = defaultYamlConfigName
	}
	cfg := &Config{}

	if err := cfg.parseYAML(configName, configPath); err != nil {
		return nil, err
	}

	if err := cfg.parseEnv(defaultLocalEnvConfigName, defaultLocalEnvConfigPath); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) parseYAML(fileName, filePath string) error {

	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filePath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return err
	}

	return nil
}

func (c *Config) parseEnv(fileName, filePath string) error {

	viper.SetConfigName(fileName)
	viper.SetConfigType("env")
	viper.AddConfigPath(filePath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&c.Postgres); err != nil {
		return err
	}

	if err := viper.Unmarshal(&c.RabbitMQ); err != nil {
		return err
	}

	if err := viper.Unmarshal(&c.Server); err != nil {
		return err
	}

	return nil
}
