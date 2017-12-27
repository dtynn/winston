package goleveldb

import (
	"os"
	"path/filepath"

	"github.com/dtynn/winston/pkg/storage"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Open return a goleveldb storage
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
		option: &opt.Options{},

		sropt: &opt.ReadOptions{},
		swopt: &opt.WriteOptions{},

		itopt: &opt.ReadOptions{},
		bwopt: &opt.WriteOptions{},
	}

	for _, o := range opts {
		o(s)
	}

	db, err := leveldb.OpenFile(path, s.option)
	if err != nil {
		return nil, err
	}

	s.db = db
	return s, nil
}

// Storage storage implementation
type Storage struct {
	dir  string
	path string

	db     *leveldb.DB
	option *opt.Options

	sropt *opt.ReadOptions
	swopt *opt.WriteOptions

	itopt *opt.ReadOptions
	bwopt *opt.WriteOptions
}

// Get return value for specified key, return nil if key not found
func (s *Storage) Get(key []byte) ([]byte, error) {
	val, err := s.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}

	return val, err
}

// MGet return values fro multiple keys
func (s *Storage) MGet(keys ...[]byte) ([][]byte, error) {
	ss, err := s.db.GetSnapshot()
	if err != nil {
		return nil, err
	}

	defer ss.Release()

	vals := make([][]byte, len(keys))
	for i, k := range keys {
		val, err := ss.Get(k, nil)
		if err == leveldb.ErrNotFound {
			continue
		}

		if err != nil {
			return nil, err
		}

		vals[i] = val
	}

	return vals, nil
}

// Put udpate the key with val
func (s *Storage) Put(key, val []byte) error {
	return s.db.Put(key, val, nil)
}

// Del delete the key
func (s *Storage) Del(key []byte) error {
	return s.db.Delete(key, nil)
}

// PrefixIterator return a iterator with prefix
func (s *Storage) PrefixIterator(prefix []byte) (storage.Iterator, error) {
	var slice *util.Range
	if prefix != nil {
		slice = util.BytesPrefix(prefix)
	}

	return s.iterator(slice)
}

// RangeIterator return a iterator within the range
func (s *Storage) RangeIterator(start, end []byte) (storage.Iterator, error) {
	var slice *util.Range
	if start != nil || end != nil {
		slice = &util.Range{
			Start: start,
			Limit: end,
		}
	}

	return s.iterator(slice)
}

func (s *Storage) iterator(slice *util.Range) (*Iterator, error) {

	return &Iterator{
		iter: s.db.NewIterator(slice, s.itopt),
	}, nil
}

// Batch open a batch
func (s *Storage) Batch() (storage.Batch, error) {
	return &Batch{
		db:    s.db,
		batch: new(leveldb.Batch),
	}, nil
}

// Close close the storage
func (s *Storage) Close() error {
	return s.db.Close()
}

// GC garbage collection
func (s *Storage) GC() error {
	return s.db.CompactRange(util.Range{})
}

func (s *Storage) cleanup() error {
	return os.RemoveAll(s.dir)
}
