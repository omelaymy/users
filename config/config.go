package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	BaseAdmin struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	Logger struct {
		ServiceName string `json:"serviceName"`
	}

	Server struct {
		Address string `json:"address"`
	}

	System struct {
		DefaultLocale string `json:"defaultLocale"`
	}
}

func LoadConfig(configFile string) (*viper.Viper, error) {
	vconfig := viper.New()
	vconfig.SetConfigFile(configFile)

	err := vconfig.ReadInConfig()
	if err != nil {
		log.Err(err).Stack().Msg("")

		return nil, fmt.Errorf("read config error: %w", err)
	}

	return vconfig, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshal config error: %w", err)
	}

	return &config, nil
}
