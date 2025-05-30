package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")

	// Try direct connection first
	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, password, host, dbName)
	db, err := sql.Open("postgres", dsn)
	if err == nil {
		// Test the connection
		err = db.Ping()
		if err == nil {
			// Connection successful, create tables
			if err := createTables(db); err != nil {
				db.Close()
				return nil, fmt.Errorf("failed to initialize database: %v", err)
			}
			return db, nil
		}
		db.Close()
	}

	// If direct connection failed, try to create everything
	// Connect to default postgres database
	defaultDSN := fmt.Sprintf("postgres://%s:%s@%s:5432/postgres?sslmode=disable", user, password, host)
	defaultDB, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		// If that fails too, try with default postgres user
		defaultDSN = fmt.Sprintf("postgres://postgres:postgres@%s:5432/postgres?sslmode=disable", host)
		defaultDB, err = sql.Open("postgres", defaultDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to postgres: %v", err)
		}
	}

	// Create database if it doesn't exist
	_, err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		defaultDB.Close()
		return nil, fmt.Errorf("failed to create database: %v", err)
	}

	defaultDB.Close()

	// Now try connecting to our database again
	db, err = sql.Open("postgres", dsn)
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
