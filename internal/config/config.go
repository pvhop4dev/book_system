package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

var config Config

func init() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("read config", "err", err)
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("unmarshal config", "err", err)
		panic(err)
	}

	slog.Info("read config", "config", config)
}

func GetConfig() Config {
	return config
}

type Config struct {
	Port string `mapstructure:"port"`
	Grpc struct {
		Port string `mapstructure:"port"`
	}
	Database struct {
		Mysql struct {
			Host     string `mapstructure:"host"`
			Port     string `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			Database string `mapstructure:"database"`
		}
	}
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	}
	Minio struct {
		Host          string `mapstructure:"host"`
		Port          string `mapstructure:"port"`
		AccessKey     string `mapstructure:"access-key"`
		SecretKey     string `mapstructure:"secret-key"`
		Location      string `mapstructure:"location"`
		DefaultBucket string `mapstructure:"default-bucket"`
		Secure        bool   `mapstructure:"secure"`
		ReturnUrl     string `mapstructure:"return-url"`
	}
}
