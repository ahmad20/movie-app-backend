package configs

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Data struct {
		Movies string
	}
	JWT struct {
		SecretKey string
		ExpiresIn time.Duration
	}
	Cors struct {
		AllowedOrigins []string
		AllowedMethods []string
		AllowedHeaders []string
	}
}

// LoadConfig loads the configuration from file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
