package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", user, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := createTables(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS qr_codes (
		id VARCHAR(255) PRIMARY KEY,
		url TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		expires_at TIMESTAMP,
		image_base64 TEXT,
		scan_count INTEGER DEFAULT 0
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	return nil
}
