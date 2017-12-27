package goleveldb

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// Iterator common iterator
type Iterator struct {
	iter iterator.Iterator
}

// First move to the first entry
func (i *Iterator) First() {
	i.iter.First()
}

// Last move to the last entry
func (i *Iterator) Last() {
	i.iter.Last()
}

// Seek move to the key equal or greater than seek. If no key exists, return false
func (i *Iterator) Seek(seek []byte) {
	i.iter.Seek(seek)
}

// Next move to the next key
func (i *Iterator) Next() bool {
	return i.iter.Next()
}

// Prev move to the previous key
func (i *Iterator) Prev() bool {
	return i.iter.Prev()
}

// Key current key of the cursor
func (i *Iterator) Key() []byte {
	key := i.iter.Key()
	if key != nil {
		res := make([]byte, len(key))
		copy(res, key)
		return res
	}

	return nil
}

// Value current value of the cursor
func (i *Iterator) Value() []byte {
	val := i.iter.Value()
	if val != nil {
		res := make([]byte, len(val))
		copy(res, val)
		return res
	}

	return nil
}

// Valid if the current entry is valid
func (i *Iterator) Valid() bool {
	return i.iter.Valid()
}

// Close close the iter
func (i *Iterator) Close() error {
	i.iter.Release()
	return i.iter.Error()
}

// Err return error if any during cursor moves
func (i *Iterator) Err() error {
	return i.iter.Error()
}
