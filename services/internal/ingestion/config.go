package ingestion

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const ConfigPrefix = "INGESTION"

type Config struct {
	BindAddress string      `envconfig:"BIND_ADDRESS" required:"true"`
	Kafka       KafkaConfig `envconfig:"KAFKA"`
}

type KafkaConfig struct {
	Brokers []string `envconfig:"BROKERS" required:"true"`
	Topic   string   `envconfig:"TOPIC" required:"true"`
	GroupID string   `envconfig:"GROUP_ID" required:"true"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := envconfig.Process(ConfigPrefix, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
