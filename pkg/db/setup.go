package db

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/safayildirim/asset-management-service/pkg/config"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"runtime"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
)

func NewConnection(conf config.PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		conf.Host, conf.User, conf.Pass, conf.DBName, conf.Port, conf.SslMode,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get raw database connection: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	err = runMigrations(conf)
	if err != nil {
		panic(err)
	}

	log.Println("Database connection verified")

	return db, nil
}

func runMigrations(conf config.PostgresConfig) error {
	_, callerDir, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Error generating env dir")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.User, conf.Pass, conf.Host, conf.Port,
		conf.DBName)

	dir := filepath.Join("file://", filepath.Dir(callerDir), "../..", "/db/migrations")

	// Initialize the migration source
	m, err := migrate.New(
		dir, // Path to migrations directory
		dsn, // PostgreSQL connection string
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")

	return nil
}
