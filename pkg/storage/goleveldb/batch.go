package goleveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Batch batch operation
type Batch struct {
	db    *leveldb.DB
	batch *leveldb.Batch
}

// Put update a key
func (b *Batch) Put(key, val []byte) error {
	b.batch.Put(key, val)
	return nil
}

// Del delete a key
func (b *Batch) Del(key []byte) error {
	b.batch.Delete(key)
	return nil
}

// Commit commit the changes
func (b *Batch) Commit() error {
	return b.db.Write(b.batch, nil)
}

// Close close the batch
func (b *Batch) Close() error {
	return nil
}
