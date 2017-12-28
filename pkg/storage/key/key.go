package key

import (
	"bytes"
	"encoding"
	"errors"
	"fmt"

	exerr "github.com/pkg/errors"
)

var (
	// ErrMalformedKeyPrefix key with malformed prefix
	ErrMalformedKeyPrefix = errors.New("malformed key prefix")
)

// Formatter the content part of a storage key
type Formatter interface {
	Bytes() []byte
	// SizeFixed bianry size fixed
}

// FixedFormatter formarter with fixed size
type FixedFormatter interface {
	Formatter
	encoding.BinaryUnmarshaler
	Size() int
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

// FixedFormatters a slice of formaters
type FixedFormatters []FixedFormatter

// Bytes implement Formater
func (f FixedFormatters) Bytes() []byte {
	buf := new(bytes.Buffer)
	for i := range f {
		buf.Write(f[i].Bytes())
	}

	return buf.Bytes()
}

// UnmarshalBinary implement encoding.BinaryUnmarshaler
func (f FixedFormatters) UnmarshalBinary(buf []byte) error {
	total := len(buf)
	read := 0

	for i, fm := range f {
		required := fm.Size()
		if len(buf) < required {
			break
		}

		if err := fm.UnmarshalBinary(buf[:required]); err != nil {
			return exerr.WithMessage(err, fmt.Sprintf("unmarshaling %T at %d", fm, i))
		}

		buf = buf[required:]
		read += required
	}

	if read != total {
		return ErrMalformedKeySize
	}

	return nil
}
