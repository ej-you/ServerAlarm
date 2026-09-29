package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnDB(t *testing.T) {
	dsn := os.Getenv("TEST_DSN")
	require.NotZero(t, dsn, "TEST_DSN env var is not specified")

	db, err := New(dsn)
	if err != nil {
		t.Fatalf("conn db: %v", err)
	}
	defer db.Close()

	t.Logf("Conn to db: %#v", db.Stats())
}
