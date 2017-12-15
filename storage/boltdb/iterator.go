package boltdb

import (
	"github.com/coreos/bbolt"
)

// Iterator common iterator
type Iterator struct {
	tx  *bolt.Tx
	cur *bolt.Cursor

	key   []byte
	val   []byte
	valid bool

	moved bool
}

// First move to the first entry
func (i *Iterator) First() {
	i.moved = true
	i.key, i.val = i.cur.First()
	i.valid = i.key != nil
}

// Seek move to the key equal or greater than seek. If no key exists, return false
func (i *Iterator) Seek(seek []byte) {
	i.moved = true
	i.key, i.val = i.cur.Seek(seek)
	i.valid = i.key != nil
}

// Next move to the next key
func (i *Iterator) Next() bool {
	if !i.moved {
		i.First()
		return i.valid
	}

	i.key, i.val = i.cur.Next()
	i.valid = i.key != nil
	return i.valid
}

// Key current key of the cursor
func (i *Iterator) Key() []byte {
	if !i.valid {
		return nil
	}

	return i.key
}

// Value current value of the cursor
func (i *Iterator) Value() []byte {
	if !i.valid {
		return nil
	}

	return i.val
}

// Valid if the current entry is valid
func (i *Iterator) Valid() bool {
	return i.valid
}

// UpdateValid for iter wrappers
func (i *Iterator) UpdateValid(valid bool) {
	i.valid = valid
}

// Close close the iter
func (i *Iterator) Close() error {
	return i.tx.Rollback()
}

// Err return error if any during cursor moves
func (i *Iterator) Err() error {
	return nil
}
