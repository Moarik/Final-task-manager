package postgre

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
	"taskManager/task/internal/adapter/postgres/dao"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
)

type Config struct {
	User         string `env:"DB_USER"`
	Password     string `env:"DB_PASSWORD"`
	Host         string `env:"DB_HOST"`
	Port         string `env:"DB_PORT"`
	DriverName   string `env:"DB_DRIVER"`
	DatabaseName string `env:"DB_NAME"`
	SSLMode      string `env:"DB_SSLMODE"`
}

type DB struct {
	Conn *gorm.DB
}

func NewDB(ctx context.Context, cfg Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	return &DB{Conn: db}, nil
}

func (db *DB) GetDB() *gorm.DB {
	return db.Conn
}

func (db *DB) Migrate() {
	if err := db.GetDB().Migrator().AutoMigrate(&dao.Task{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
