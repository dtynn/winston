package boltdb

import (
	"os"
	"path/filepath"

	"github.com/coreos/bbolt"
	"github.com/dtynn/winston/pkg/storage"
)

// Open return a boltdb storage
func Open(path string, opts ...Option) (*Storage, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(abs)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	s := &Storage{
		dir:    dir,
		path:   abs,
		bucket: defaultBucket,
		opt:    *bolt.DefaultOptions,
	}

	for _, o := range opts {
		o(s)
	}

	db, err := bolt.Open(path, 0600, &s.opt)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(s.bucket)
		return err

	}); err != nil {
		return nil, err
	}

	s.db = db
	return s, nil
}

// Storage storage implementation
type Storage struct {
	dir    string
	path   string
	bucket []byte

	opt bolt.Options
	db  *bolt.DB
}

// Get return value for specified key, return nil if key not found
func (s *Storage) Get(key []byte) ([]byte, error) {
	var val []byte

	if err := s.db.View(func(tx *bolt.Tx) error {
		val = tx.Bucket(s.bucket).Get(key)
		return nil
	}); err != nil {
		return nil, err
	}

	return val, nil
}

// MGet return values fro multiple keys
func (s *Storage) MGet(keys ...[]byte) ([][]byte, error) {
	vals := make([][]byte, len(keys))

	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		for i, key := range keys {
			vals[i] = b.Get(key)
		}

		return nil

	}); err != nil {
		return nil, err
	}

	return vals, nil
}

// Put udpate the key with val
func (s *Storage) Put(key, val []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(s.bucket).Put(key, val)
	})
}

// Del delete the key
func (s *Storage) Del(key []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(s.bucket).Delete(key)
	})
}

// PrefixIterator return a iterator with prefix
func (s *Storage) PrefixIterator(prefix []byte) (storage.Iterator, error) {
	return s.iterator(prefix, storage.PrefixEnd(prefix))
}

// RangeIterator return a iterator within the range
func (s *Storage) RangeIterator(start, end []byte) (storage.Iterator, error) {
	return s.iterator(start, end)
}

func (s *Storage) iterator(start, end []byte) (*Iterator, error) {
	tx, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}

	return &Iterator{
		start: start,
		end:   end,
		tx:    tx,
		cur:   tx.Bucket(s.bucket).Cursor(),
	}, nil
}

// Batch open a batch
func (s *Storage) Batch() (storage.Batch, error) {
	return &Batch{
		s: s,
		batchOp: batchOp{
			put: make([][2][]byte, 0, 100),
			del: make([][]byte, 0, 100),
		},
	}, nil
}

// Close close the storage
func (s *Storage) Close() error {
	return s.db.Close()
}

// GC garbage collection
func (s *Storage) GC() error {
	return nil
}

func (s *Storage) cleanup() error {
	return os.RemoveAll(s.dir)
}
