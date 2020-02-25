package kvstore

import (
	"io"
)

// Store is generic interface for local storages with string key
type Store interface {
	// NewEntry cretae data entry and return new key
	NewEntry(data []byte) ([]byte, error)
	// Save entry. Returns error if key already exists
	Save(key []byte, data []byte) error
	// Update or create entry.
	Update(key []byte, data []byte) error
	// Load data with key. Return nil if data is absent
	Load(key []byte) ([]byte, error)
	// Timestamp of last key update. Zero if data is absent
	Timestamp(key []byte) (int64, error)
	io.Closer
}
