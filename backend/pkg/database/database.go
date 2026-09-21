package database

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB(log *logrus.Logger) (*gorm.DB, func(), error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	appEnv := os.Getenv("APP_ENV")
	logLevel := logger.Warn
	if appEnv == "development" || appEnv == "staging" {
		logLevel = logger.Info
	}

	newLogger := logger.New(log, logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: appEnv == "production",
		Colorful:                  true,
	})

	db := connectWithRetry(dsn, newLogger, 5)

	cleanup := func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.WithError(err).Error("failed to get sql db")
			return
		}

		err = sqlDB.Close()
		if err != nil {
			log.WithError(err).Error("failed to close sql db")
			return
		}
		log.Info("sql db closed")
	}

	return db, cleanup, nil
}

func connectWithRetry(dsn string, dbLogger logger.Interface, maxRetries int) *gorm.DB {
	var gormDB *gorm.DB
	var err error

	for i := range maxRetries {
		gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:         dbLogger,
			TranslateError: true,
		})
		if err == nil {
			sqlDB, dbErr := gormDB.DB()
			if dbErr == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					configurePool(gormDB)
					log.Println("✅ Database connected")
					return gormDB
				} else {
					err = pingErr
				}
			} else {
				err = dbErr
			}
		}

		waitTime := time.Duration(i+1) * 2 * time.Second
		log.Printf("⚠️ DB connection attempt %d/%d failed: %v. Retrying in %v...",
			i+1, maxRetries, err, waitTime)
		time.Sleep(waitTime)
	}

	log.Fatalf("❌ Failed to connect to database after %d retries: %v", maxRetries, err)
	return nil
}

func configurePool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
}
