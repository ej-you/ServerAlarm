// Package storage contains storage interfaces to store data.
package storage

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found or expired") // key not found error

// KeyValue describes simple key-value storage.
type KeyValue interface {
	// Get returns value by the given key.
	// Returns error if key-value was not found.
	Get(key string) ([]byte, error)
	// Set sets new key-value pair (expired after exp).
	// Expiration value 0 means no expiration.
	Set(key string, val []byte, exp time.Duration) error
	// Delete deletes value by the given key.
	// It handles non-existent keys gracefully (without error).
	Delete(key string) error
	// Reset deletes all key-values from the storage.
	Reset() error
	// Close closes the storage.
	Close() error
}
