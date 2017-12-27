package key

import (
	"bytes"
	"encoding"
	"errors"
)

var (
	// ErrMalformedKeyPrefix key with malformed prefix
	ErrMalformedKeyPrefix = errors.New("malformed key prefix")
)

// Formatter the content part of a storage key
type Formatter interface {
	Bytes() []byte
}

// Key return a storage key of parts
func Key(prefix []byte, formater Formatter) []byte {
	f := formater.Bytes()
	key := make([]byte, len(prefix)+len(f))
	copy(key, prefix)
	copy(key[len(prefix):], f)
	return key
}

// Unmarshal get the content of the key
func Unmarshal(key, prefix []byte, v encoding.BinaryUnmarshaler) error {
	if !bytes.HasPrefix(key, prefix) {
		return ErrMalformedKeyPrefix
	}

	return v.UnmarshalBinary(key[len(prefix):])
}
