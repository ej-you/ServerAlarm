package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // DB driver
)

const _pingTimeout = 5 * time.Second // ping timeout to check DB is ready

// DB represents a DB connection.
type DB struct {
	*sql.DB
}

// New returns a new instance of DB with connection to the given DSN.
func New(dsn string) (*DB, error) {
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db := &DB{sqlDB}
	// check connection
	if err := db.IsReady(); err != nil {
		return nil, err
	}
	return db, nil
}

// IsReady checks that DB is ready to use.
// Returns nil error if DB is ready.
func (db *DB) IsReady() error {
	ctx, cancel := context.WithTimeout(context.Background(), _pingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("db is unreachable: %w", err)
	}
	return nil
}
