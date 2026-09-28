package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// New returns a new instance of DB with connection to the given DSN.
func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	// check connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db is unreachable: %w", err)
	}
	return db, nil
}
