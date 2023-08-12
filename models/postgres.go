package models

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Callers of Open need to ensure that the
// connection is eventually closed via the
// db.Close() method
func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return db, nil
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "hashem",
		Password: "4g5gbook",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func Migrate(db *sql.DB, dir string) error {
	// connect goose with our database
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	// run the upgrade command
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	if dir == "" {
		dir = "."
	}
	// set a global variable that will set the migration as an embedding
	goose.SetBaseFS(migrationFS)

	// to let other programmes to run after embedding the sql code
	// is to set it to nil after finishing
	defer func() {
		// undo the file system main set
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}
