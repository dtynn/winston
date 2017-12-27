package storage

// Storage storage interface
type Storage interface {
	// if key not found, just return nil value
	Get(key []byte) ([]byte, error)

	// if any key not found, just return nil value
	MGet(keys ...[]byte) ([][]byte, error)

	Put(key, val []byte) error
	Del(key []byte) error

	PrefixIterator(prefix []byte) (Iterator, error)
	RangeIterator(start, end []byte) (Iterator, error)

	Batch() (Batch, error)

	Close() error
	GC() error
}

// Iterator iter interface
type Iterator interface {
	First()
	Seek(seek []byte)
	Next() bool
	Key() []byte
	Value() []byte
	Valid() bool

	Close() error
	Err() error
}

// ManagedIterator used for prefix iterator wrapper or range iterator wrapper
type ManagedIterator interface {
	Iterator
	UpdateValid(valid bool)
}

// Batch batch interface
type Batch interface {
	Put(key, val []byte) error
	Del(key []byte) error

	Commit() error
	Close() error
}
