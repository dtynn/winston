package storage

import "bytes"

var (
	_ Iterator = (*prefixIterator)(nil)
)

// PrefixIterator return a wrapped iterator
func PrefixIterator(prefix []byte, iter ManagedIterator) Iterator {
	if prefix == nil {
		return iter
	}

	wrapped := &prefixIterator{
		prefix:          prefix,
		ManagedIterator: iter,
	}

	return wrapped
}

type prefixIterator struct {
	prefix []byte
	ManagedIterator
	moved bool
}

func (p *prefixIterator) checkKeyValid() bool {
	key := p.ManagedIterator.Key()

	return key != nil && bytes.HasPrefix(key, p.prefix)
}

func (p *prefixIterator) First() {
	p.Seek(p.prefix)
}

func (p *prefixIterator) Seek(seek []byte) {
	// check if the seek has prefix
	if !bytes.HasPrefix(seek, p.prefix) {
		if bytes.Compare(seek, p.prefix) < 0 {
			seek = p.prefix
		} else {
			p.ManagedIterator.UpdateValid(false)
			return
		}
	}

	p.moved = true
	p.ManagedIterator.Seek(seek)
	p.ManagedIterator.UpdateValid(p.checkKeyValid())
}

func (p *prefixIterator) Next() bool {
	if !p.moved {
		p.First()
		return p.ManagedIterator.Valid()
	}

	if ok := p.ManagedIterator.Next(); !ok {
		return false
	}

	p.ManagedIterator.UpdateValid(p.checkKeyValid())
	return p.ManagedIterator.Valid()
}
