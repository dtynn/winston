package storage

import "bytes"

var (
	_ Iterator = (*rangeIterator)(nil)
)

// RangeIterator return a wrapped iterator
func RangeIterator(start, end []byte, iter ManagedIterator) Iterator {
	if start == nil && end == nil {
		return iter
	}

	wrapped := &rangeIterator{
		start:           start,
		end:             end,
		ManagedIterator: iter,
	}

	return wrapped
}

type rangeIterator struct {
	start, end []byte
	ManagedIterator
	moved bool
}

func (r *rangeIterator) checkKeyValid() bool {
	key := r.ManagedIterator.Key()
	valid := key != nil
	if valid && r.end != nil {
		valid = bytes.Compare(key, r.end) < 0
	}

	return valid
}

func (r *rangeIterator) First() {
	r.Seek(r.start)
}

func (r *rangeIterator) Seek(seek []byte) {
	if r.start != nil && bytes.Compare(seek, r.start) < 0 {
		seek = r.start
	}

	r.moved = true

	r.ManagedIterator.Seek(seek)
	r.ManagedIterator.UpdateValid(r.checkKeyValid())
}

func (r *rangeIterator) Next() bool {
	if !r.moved {
		r.First()
		return r.ManagedIterator.Valid()
	}

	if ok := r.ManagedIterator.Next(); !ok {
		return false
	}

	r.ManagedIterator.UpdateValid(r.checkKeyValid())
	return r.ManagedIterator.Valid()
}
