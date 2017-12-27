package boltdb

import (
	"github.com/coreos/bbolt"
	"github.com/dtynn/winston/pkg/storage"
)

// Iterator common iterator
type Iterator struct {
	tx  *bolt.Tx
	cur *bolt.Cursor

	key   []byte
	val   []byte
	valid bool

	start, end []byte

	moved bool
}

// First move to the first entry
func (i *Iterator) First() {
	i.moved = true

	if i.start == nil {
		i.key, i.val = i.cur.First()
		i.valid = i.key != nil && storage.KeyInRange(i.key, nil, i.end)
		return
	}

	i.Seek(i.start)
}

// Last move to the last entry
func (i *Iterator) Last() {
	i.moved = true

	if i.end == nil {
		i.key, i.val = i.cur.Last()
		i.valid = i.key != nil && storage.KeyInRange(i.key, i.start, nil)
		return
	}

	i.cur.Seek(i.end)
	i.Prev()
}

// Seek move to the key equal or greater than seek. If no key exists, return false
func (i *Iterator) Seek(seek []byte) {
	i.moved = true

	if !storage.KeyInRange(seek, i.start, nil) {
		seek = i.start
	}

	i.key, i.val = i.cur.Seek(seek)
	i.valid = i.key != nil && storage.KeyInRange(i.key, nil, i.end)
}

// Next move to the next key
func (i *Iterator) Next() bool {
	if !i.moved {
		i.First()
		return i.valid
	}

	i.key, i.val = i.cur.Next()
	i.valid = i.key != nil && storage.KeyInRange(i.key, nil, i.end)
	return i.valid
}

// Prev move to the previous key
func (i *Iterator) Prev() bool {
	i.key, i.val = i.cur.Prev()
	i.valid = i.key != nil && storage.KeyInRange(i.key, i.start, nil)
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

// Close close the iter
func (i *Iterator) Close() error {
	return i.tx.Rollback()
}

// Err return error if any during cursor moves
func (i *Iterator) Err() error {
	return nil
}
