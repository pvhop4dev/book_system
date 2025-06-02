package main

import (
	"book_system/i18n"
	_ "book_system/internal/config"
	"log/slog"
	"os"
	// "gitlab.ai-vlab.com/cygate/common/pkg/config/minio"
	// "gitlab.ai-vlab.com/cygate/common/pkg/config/redis"
)

func init() {
	i18n.InitI18n([]string{"vi", "en"})
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})))
}

func main() {
	// redis.InitRedis()
	// minio.InitMinioClient()
	// defer redis.Close()
	slog.Info("--------------------------------------------RUNNING!--------------------------------------------")
	// router.InitRouter().Run(":" + viper.GetString("port"))
}
