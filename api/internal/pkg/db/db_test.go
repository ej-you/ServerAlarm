package db

import (
	"os"
	"testing"
)

func TestConnDB(t *testing.T) {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		t.Fatal("TEST_DSN env var is not specified")
	}

	db, err := New(dsn)
	if err != nil {
		t.Fatalf("conn db: %v", err)
	}
	defer db.Close()

	t.Logf("Conn to db: %#v", db.Stats())
}
