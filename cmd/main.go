package main

import (
	"book_system/i18n"
	"book_system/internal/config"

	"log/slog"
	"os"
	// "gitlab.ai-vlab.com/cygate/common/pkg/config/minio"
	// "gitlab.ai-vlab.com/cygate/common/pkg/config/redis"
)

func init() {
	i18n.InitI18n([]string{"vi", "en"})
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	logHandler.WithAttrs([]slog.Attr{
		slog.String("app", "book_system"),
		slog.String("version", "1.0.0"),
		slog.String("environment", config.MustGet().Environment),
	})

	// Set the default logger to use the custom handler
	slog.SetDefault(slog.New(logHandler))
	slog.Info("--------------------------------------------STARTING!--------------------------------------------")
}

func main() {
	// redis.InitRedis()
	// minio.InitMinioClient()
	// defer redis.Close()
	slog.Info("--------------------------------------------RUNNING!--------------------------------------------")

	slog.Info("Port: ", slog.Any("all", config.MustGet()))
	slog.Info("Environment: ", slog.String("environment", config.MustGet().Environment))

	// router.InitRouter().Run(":" + viper.GetString("port"))
}
