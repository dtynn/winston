package boltdb

import (
	"os"
	"path/filepath"

	"github.com/coreos/bbolt"
	"github.com/dtynn/winston/storage"
)

var (
	defaultBucket = []byte("_winston")
)

// Option db option
type Option func(s *Storage)

// Bucket modify bucket name
func Bucket(name []byte) Option {
	return func(s *Storage) {
		if len(name) > 0 {
			s.bucket = name
		}
	}
}

// Open return a boltdb storage
func Open(path string, opts ...Option) (*Storage, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	s := &Storage{
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
	iter, err := s.iterator()
	if err != nil {
		return nil, err
	}

	if prefix != nil {
		return storage.PrefixIterator(prefix, iter), nil
	}

	return iter, nil
}

// RangeIterator return a iterator within the range
func (s *Storage) RangeIterator(start, end []byte) (storage.Iterator, error) {
	iter, err := s.iterator()
	if err != nil {
		return nil, err
	}

	if start != nil || end != nil {
		return storage.RangeIterator(start, end, iter), nil
	}

	return iter, nil
}

func (s *Storage) iterator() (*Iterator, error) {
	tx, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}

	return &Iterator{
		tx:  tx,
		cur: tx.Bucket(s.bucket).Cursor(),
	}, nil
}

// Batch open a batch
func (s *Storage) Batch() (storage.Batch, error) {
	tx, err := s.db.Begin(true)
	if err != nil {
		return nil, err
	}

	return &Batch{
		tx:     tx,
		bucket: tx.Bucket(s.bucket),
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
	return os.RemoveAll(s.path)
}
