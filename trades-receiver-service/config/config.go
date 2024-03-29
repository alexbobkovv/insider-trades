package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		App        `mapstructure:"app"`
		HTTPServer `mapstructure:"http_server"`
		GRPCServer `mapstructure:"grpc_server"`
		Postgres   `mapstructure:"postgres"`
		Logger     `mapstructure:"logger"`
		RabbitMQ   `mapstructure:"rabbitmq"`
	}

	App struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
	}

	HTTPServer struct {
		Port         string `mapstructure:"port"`
		ReceiverPath string `mapstructure:"RECEIVER_PATH"`
		AllowOrigin  string `mapstructure:"allow_origin"`
	}

	GRPCServer struct {
		Port string `mapstructure:"port"`
	}

	Postgres struct {
		URL string `mapstructure:"POSTGRES_URL"`
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

	viper.AddConfigPath(filePath)
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")

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

	if err := viper.Unmarshal(&c.HTTPServer); err != nil {
		return err
	}

	return nil
}
