package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Adejare77/blogPlatform/internals/schemas"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	User     string
	Password string
	DB       string
	Host     string
	Port     int
}

var DB *gorm.DB

func Connect() error {
	cfg, err := loadDBConfig()
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	// Data source name
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Taipei",
		cfg.User, cfg.Password, cfg.Host, cfg.DB, cfg.Port)

	// Configure Database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("(DB Configuration) %s", err)
	}

	// Set DB settings
	SQLDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("(DB Settings) %s", err)
	}

	connMaxLifetime, _ := time.ParseDuration(os.Getenv("CONN_MAX_LIFETIME"))
	maxIdleConns, _ := strconv.Atoi(os.Getenv("MAX_IDLE_CONNS"))
	maxOpenConns, _ := strconv.Atoi(os.Getenv("MAX_OPEN_CONNS"))

	SQLDB.SetConnMaxLifetime(connMaxLifetime)
	SQLDB.SetMaxIdleConns(maxIdleConns)
	SQLDB.SetMaxOpenConns(maxOpenConns)

	// Auto-Migrate tables Created
	err = db.AutoMigrate(&schemas.User{}, &schemas.Post{}, &schemas.Comment{}, &schemas.Like{})
	if err != nil {
		return fmt.Errorf("(DB AutoMigration) %s", err)
	}

	DB = db

	return nil
}

func loadDBConfig() (*DBConfig, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("(DB_PORT) %s ", err)
	}
	return &DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DB:       os.Getenv("DB_DATABASE"),
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
	}, nil
}
