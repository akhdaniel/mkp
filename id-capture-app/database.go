package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")

	// Create table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS id_documents (
		id SERIAL PRIMARY KEY,
		document_id VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255),
		document_number VARCHAR(255),
		birth_date DATE,
		issue_date DATE,
		expiry_date DATE,
		image_path VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	log.Println("Database table checked/created successfully")
	return db, nil
}