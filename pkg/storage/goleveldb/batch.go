package goleveldb

import (
	"sync"

	"github.com/dtynn/winston/pkg/storage"
	"github.com/syndtr/goleveldb/leveldb"
)

// Batch batch operation
type Batch struct {
	db    *leveldb.DB
	batch *leveldb.Batch
	sync.Mutex
	closed bool
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
	b.Lock()
	defer b.Unlock()

	if b.closed {
		return storage.ErrBatchClosed
	}

	if err := b.db.Write(b.batch, nil); err != nil {
		return err
	}

	b.closed = true

	return nil
}

// Close close the batch
func (b *Batch) Close() error {
	b.Lock()
	defer b.Unlock()

	if b.closed {
		return storage.ErrBatchClosed
	}

	b.closed = true

	return nil
}
