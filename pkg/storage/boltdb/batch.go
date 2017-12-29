package boltdb

import (
	"sync"

	"github.com/coreos/bbolt"
	"github.com/dtynn/winston/pkg/storage"
)

type batchOp struct {
	put [][2][]byte
	del [][]byte
}

// Batch batch operation
type Batch struct {
	s *Storage
	batchOp
	closed bool
	sync.Mutex
}

// Put update a key
func (b *Batch) Put(key, val []byte) error {
	b.batchOp.put = append(b.batchOp.put, [2][]byte{key, val})
	return nil
}

// Del delete a key
func (b *Batch) Del(key []byte) error {
	b.batchOp.del = append(b.batchOp.del, key)
	return nil
}

// Commit commit the changes
func (b *Batch) Commit() error {
	b.Lock()
	defer b.Unlock()

	if b.closed {
		return storage.ErrBatchClosed
	}

	if err := b.s.db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.s.bucket)

		for _, pair := range b.batchOp.put {
			if err := bucket.Put(pair[0], pair[1]); err != nil {
				return err
			}
		}

		for _, key := range b.batchOp.del {
			if err := bucket.Delete(key); err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
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
