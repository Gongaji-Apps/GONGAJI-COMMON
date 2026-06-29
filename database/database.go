package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host         string
	User         string
	Password     string
	DBName       string
	Port         string
	SSLMode      string
	TimeZone     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	LogLevel     logger.LogLevel

	// StatementTimeout (opsional) menyetel server-side statement_timeout per koneksi
	// sebagai pengaman query lambat. 0 = tidak diaktifkan (perilaku lama).
	StatementTimeout time.Duration
}

func NewPostgres(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	if cfg.StatementTimeout > 0 {
		dsn += fmt.Sprintf(" options='-c statement_timeout=%d'", cfg.StatementTimeout.Milliseconds())
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	return db, nil
}
