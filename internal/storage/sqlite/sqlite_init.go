package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type UserLicense struct {
	ID        int64
	License   string
	UserId    string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	HWID      *string
	Status    string
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	// Open the SQLite database
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Prepare SQL statements to create the required tables if they don't exist
	createUserLicenseTable := `
CREATE TABLE IF NOT EXISTS UserLicense (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    license VARCHAR(25) NOT NULL UNIQUE,
    UserId VARCHAR(50) NOT NULL UNIQUE,
    createdAt TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP NOT NULL,
    expiresAt TIMESTAMP NOT NULL,
    hwid VARCHAR(50),
    status VARCHAR(10) NOT NULL
);`

	createTransactionLogsTable := `
CREATE TABLE IF NOT EXISTS TransactionLogs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    description TEXT NOT NULL
);`

	// Execute the SQL statement to create UserLicense table
	if _, err := db.Exec(createUserLicenseTable); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Execute the SQL statement to create TransactionLogs table
	if _, err := db.Exec(createTransactionLogsTable); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Return the storage instance
	return &Storage{db: db}, nil
}
