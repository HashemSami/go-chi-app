package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
		"postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func main() {
	// this doesn't actually communicating with the server
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "hashem",
		Password: "4g5gbook",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
	fmt.Println(cfg.String())
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbpool, err := pgxpool.New(context.Background(), cfg.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

	createTables(dbpool)
}

func createTables(dbpool *pgxpool.Pool) {
	_, err := dbpool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		age INT,
		first_name TEXT,
		last_name TEXT,
		email TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS orders(
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		amount INT,
		descriptions TEXT
	);
	`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tables created.")
}
