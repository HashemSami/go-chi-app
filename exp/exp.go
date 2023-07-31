package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
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

	// createTables(db)
	// insertData(db)
	// insertRow(db)
	// queryUser(db)
	// createFakeOrders(db)
	queryOrders(db)
}

func createTables(db *sql.DB) {
	_, err := db.Exec(`
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
		description TEXT
	);
	`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tables created.")
}

func insertData(db *sql.DB) {
	name := "Hashem5"
	email := "Hash5@hash.com"

	_, err := db.Exec(
		`INSERT INTO users(first_name, email)
	VALUES($1, $2)`,
		name, email)
	if err != nil {
		panic(err)
	}

	fmt.Println("Data inserted.")
}

func insertRow(db *sql.DB) {
	name := "Hashem7"
	email := "Hash7@hash.com"

	row := db.QueryRow(
		`INSERT INTO users(first_name, email)
	VALUES($1, $2) RETURNING id`,
		name, email)

	var id int
	err := row.Scan(&id)
	if err != nil {
		panic(err)
	}

	fmt.Println("row created with id:", id)
}

func queryUser(db *sql.DB) {
	email := "Hash2@hash.com"
	row := db.QueryRow(
		`SELECT first_name from users where email=$1`, email)

	var name string
	err := row.Scan(&name)
	if err == sql.ErrNoRows {
		fmt.Println("No Row Found!!!")
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("name retirved:", name)
}

func createFakeOrders(db *sql.DB) {
	userId := 1
	for i := 1; i <= 5; i++ {
		amount := i * 100
		desc := fmt.Sprintf("Fake order #%d", i)

		_, err := db.Exec(`
		INSERT INTO orders(user_id, amount, description)
		VALUES($1,$2,$3)`, userId, amount, desc)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("added fake orders.")
}

func queryOrders(db *sql.DB) {
	type Order struct {
		ID     int
		UserID int
		Amount int
		Desc   string
	}

	var orders []Order

	userID := 1

	rows, err := db.Query(`
	SELECT id, amount, description
	FROM orders
	WHERE user_id=$1`, userID)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	// rows will be empty until you call the first next function
	// on it
	for rows.Next() {
		var order Order
		order.UserID = userID
		err := rows.Scan(&order.ID, &order.Amount, &order.Desc)
		if err != nil {
			panic(err)
		}
		orders = append(orders, order)
	}
	// check for error happened on the loop
	if rows.Err() != nil {
		panic(rows.Err())
	}

	fmt.Println("Orders:", orders)
}
