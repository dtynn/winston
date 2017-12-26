package key

import (
	"encoding/binary"
	"errors"
)

var (
	// ErrMalformedKeySize malformed buf size for a key
	ErrMalformedKeySize = errors.New("malformed key size")
)

// UI64 key formater wrapper of uint64
type UI64 uint64

// Bytes return bytes
func (u UI64) Bytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(u))
	return buf
}

// UnmarshalBinary get uint64 from bytes
func (u *UI64) UnmarshalBinary(buf []byte) error {
	if len(buf) != 8 {
		return ErrMalformedKeySize
	}

	*u = UI64(binary.BigEndian.Uint64(buf))
	return nil
}

// I64 key formater wrapper of int64
type I64 int64

// Bytes return bytes
func (i I64) Bytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// UnmarshalBinary get int64 from bytes
func (i *I64) UnmarshalBinary(buf []byte) error {
	if len(buf) != 8 {
		return ErrMalformedKeySize
	}

	*i = I64(binary.BigEndian.Uint64(buf))
	return nil
}

// UI32 key formater wrapper of uint32
type UI32 uint32

// Bytes return bytes
func (u UI32) Bytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(u))
	return buf
}

// UnmarshalBinary get uint32 from bytes
func (u *UI32) UnmarshalBinary(buf []byte) error {
	if len(buf) != 4 {
		return ErrMalformedKeySize
	}

	*u = UI32(binary.BigEndian.Uint32(buf))
	return nil
}

// I32 key formater wrapper of int32
type I32 int32

// Bytes return bytes
func (i I32) Bytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

// UnmarshalBinary get int32 from bytes
func (i *I32) UnmarshalBinary(buf []byte) error {
	if len(buf) != 4 {
		return ErrMalformedKeySize
	}

	*i = I32(binary.BigEndian.Uint32(buf))
	return nil
}
