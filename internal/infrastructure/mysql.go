package infrastructure

import (
	"book_system/internal/config"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
	dbErr  error
)

// InitDB initializes the database connection pool
func InitDB() (*gorm.DB, error) {
	dbOnce.Do(func() {
		cfg := config.MustGet()
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.Mysql.User,
			cfg.Database.Mysql.Password,
			cfg.Database.Mysql.Host,
			cfg.Database.Mysql.Port,
			cfg.Database.Mysql.Database,
		)

		slog.Info("Connecting to MySQL...", "dsn", dsn)

		gormConfig := &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // Change to logger.Silent in production
		}

		db, dbErr = gorm.Open(mysql.Open(dsn), gormConfig)
		if dbErr != nil {
			slog.Error("Failed to connect to database", "error", dbErr)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			dbErr = err
			slog.Error("Failed to get database instance", "error", err)
			return
		}

		// Set connection pool parameters
		sqlDB.SetMaxIdleConns(10)                  // Maximum number of idle connections
		sqlDB.SetMaxOpenConns(100)                 // Maximum number of open connections
		sqlDB.SetConnMaxLifetime(time.Hour)        // Maximum connection lifetime
		sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Maximum idle time

		slog.Info("Successfully connected to MySQL database")
	})

	return db, dbErr
}

// GetDB returns the database instance
// Panics if the database is not initialized
func GetDB() *gorm.DB {
	if db == nil {
		slog.Error("Database not initialized, call InitDB() first")
		panic("database not initialized")
	}
	return db
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			slog.Error("Failed to get database instance for closing", "error", err)
			return
		}
		sqlDB.Close()
		slog.Info("Closed database connection")
	}
}
