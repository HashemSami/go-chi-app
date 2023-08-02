package main

import (
"database/sql"
"fmt"

    _ "github.com/jackc/pgx/v5/stdlib"

)

type PostgresConfig struct {
Host string
Port string
User string
Password string
Database string
SSLMode string
}

func (cfg PostgresConfig) String() string {
return fmt.Sprintf(
"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func main() {
// this doesn't actually communicating with the server
cfg := PostgresConfig{
Host: "localhost",
Port: "5432",
User: "hashem",
Password: "4g5gbook",
Database: "lenslocked",
SSLMode: "disable",
}
fmt.Println(cfg.String())

    db, err := sql.Open("pgx", cfg.String())
    if err != nil {
    	panic(err)
    }

    defer db.Close()

    // to make sure that the database is responding
    err = db.Ping()

    if err != nil {
    	panic(err)
    }

    fmt.Println("connected...")

}

====================================================================
package main

import (
"context"
"fmt"
"os"

    "github.com/jackc/pgx/v5/pgxpool"

)

type PostgresConfig struct {
Host string
Port string
User string
Password string
Database string
SSLMode string
}

func (cfg PostgresConfig) String() string {
return fmt.Sprintf(
"postgres://%s:%s@%s:%s/%s",
cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func main() {
// this doesn't actually communicating with the server
cfg := PostgresConfig{
Host: "localhost",
Port: "5432",
User: "hashem",
Password: "4g5gbook",
Database: "lenslocked",
SSLMode: "disable",
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

}
