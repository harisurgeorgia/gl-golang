package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var Conn *sql.DB

func Init() {
	var err error
	dsn := "host=localhost port=5432 user=postgres password=postgress dbname=postgres sslmode=disable"
	Conn, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	if err = Conn.Ping(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
}
