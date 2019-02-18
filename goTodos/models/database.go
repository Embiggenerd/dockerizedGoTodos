package models

import (
	"database/sql"
	"fmt" // Package pq importing drivers for db
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	dbname   = "go"
	password = "postgres"
)

var db *sql.DB

// Init initializes our db in main
func Init() {
	var err error

	psqlInfo := os.Getenv("DATABASE_URL")

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	createTables()

	fmt.Println("Successfully connected!")
}

func createTables() {
	_, err := db.Query(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			age INT,
			first_name TEXT,
			last_name TEXT,
			email TEXT UNIQUE NOT NULL,
			password TEXT
			);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Query(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			body TEXT,
			authorid INT,
			done BOOLEAN
			);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Query(`
		CREATE TABLE IF NOT EXISTS sessions (
			id SERIAL PRIMARY KEY,
			userid INT,
			hex TEXT
			);`)
	if err != nil {
		panic(err)
	}
}
