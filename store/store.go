package store

import (
	"io"
)

// Store is generic interface for local storages with string key
type Store interface {
	// Save data with key
	Save(key []byte, data []byte) error
	// Load data from key. Return nil if data is absent
	Load(key []byte) ([]byte, error)
	// Timestamp of last key update. Zero if data is absent
	Timestamp(key []byte) (int64, error)
	io.Closer
}

// KeyStore allows create new entires with random unique keys
type KeyStore interface {
	Store
	// NewKey data and return new key
	NewKey(data []byte) ([]byte, error)
}
