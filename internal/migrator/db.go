package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dsn string) {
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
}

func CreateMigrationsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		status TEXT NOT NULL,
		applied_at TIMESTAMP
	)`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations table ensured.")
}
