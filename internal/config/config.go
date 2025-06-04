package config

import (
	"log/slog"
	"sync"

	"github.com/spf13/viper"
)

type config struct {
	Environment string `mapstructure:"environment"`
	Port        uint32 `mapstructure:"port"`
	Grpc        struct {
		Port string `mapstructure:"port"`
	}
	JWT struct {
		AccessSecret  string `mapstructure:"access-secret"`
		RefreshSecret string `mapstructure:"refresh-secret"`
		AccessExpiry  int    `mapstructure:"access-expiry"`
		RefreshExpiry int    `mapstructure:"refresh-expiry"`
	}
	RateLimiter struct {
		Burst int `mapstructure:"burst"`
		Rate  int `mapstructure:"rate"`
	}
	Codec struct {
		SecretKey uint32 `mapstructure:"secret-key"`
	}
	Casbin struct {
		DSN string `mapstructure:"dsn"`
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
		ReturnURL     string `mapstructure:"return-url"`
	}
}

var (
	once     sync.Once
	instance *config
)

func setDefault() {
	viper.SetDefault("port", "8888")
	viper.SetDefault("database.mysql.host", "default")
	viper.SetDefault("database.mysql.port", "default")
	viper.SetDefault("database.mysql.user", "default")
	viper.SetDefault("database.mysql.password", "default")
	viper.SetDefault("database.mysql.database", "default")
	viper.SetDefault("grpc.port", "default")
}

// Initialize loads and validates the configuration.
// It's safe to call this function multiple times, but the configuration
// will only be loaded once.
func init() {
	once.Do(func() {
		setDefault()
		viper.SetConfigFile("config.yaml")
		if err := viper.ReadInConfig(); err != nil {
			slog.Error("failed to read config file", "error", err)
			panic(err)
		}

		var cfg config
		if err := viper.Unmarshal(&cfg); err != nil {
			slog.Error("failed to unmarshal config", "error", err)
			panic(err)
		}

		instance = &cfg
		slog.Info("configuration loaded successfully")
	})
}

func get() *config {
	if instance == nil {
		panic("config not initialized, call Initialize() first")
	}
	return instance
}

// MustGet is like Get but returns a non-pointer value.
// It's useful for making the configuration read-only.
func MustGet() config {
	return *get()
}
