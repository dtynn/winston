package boltdb

import (
	"github.com/coreos/bbolt"
)

// Batch batch operation
type Batch struct {
	tx     *bolt.Tx
	bucket *bolt.Bucket
}

// Put update a key
func (b *Batch) Put(key, val []byte) error {
	return b.bucket.Put(key, val)
}

// Del delete a key
func (b *Batch) Del(key []byte) error {
	return b.bucket.Delete(key)
}

// Commit commit the changes
func (b *Batch) Commit() error {
	return b.tx.Commit()
}

// Close close the batch
func (b *Batch) Close() error {
	return b.tx.Rollback()
}
