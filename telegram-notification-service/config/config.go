package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		App    `mapstructure:"app"`
		Server `mapstructure:"server"`
		Logger `mapstructure:"logger"`
		Telegram
		RabbitMQ `mapstructure:"rabbitmq"`
	}

	App struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
	}

	Telegram struct {
		BotToken  string `mapstructure:"TELEGRAM_BOT_TOKEN"`
		ChannelID int64  `mapstructure:"TELEGRAM_CHANNEL_ID"`
	}

	Server struct {
		Port string `mapstructure:"port"`
	}

	Logger struct {
		Level    string `mapstructure:"level"`
		Format   string `mapstructure:"format"`
		Filepath string `mapstructure:"filepath"`
	}

	RabbitMQ struct {
		AmqpURL      string `mapstructure:"AMQP_URL"`
		Exchange     string `mapstructure:"exchange"`
		Durable      bool   `mapstructure:"durable"`
		QueueName    string `mapstructure:"queueName"`
		RoutingKey   string `mapstructure:"routingKey"`
		ConsumerName string `mapstructure:"consumerName"`
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

	if err := viper.Unmarshal(&c.Telegram); err != nil {
		return err
	}

	if err := viper.Unmarshal(&c.RabbitMQ); err != nil {
		return err
	}

	return nil
}
