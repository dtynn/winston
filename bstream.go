package winston

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"io"
)

var (
	_ encoding.BinaryMarshaler   = (*bstream)(nil)
	_ encoding.BinaryUnmarshaler = (*bstream)(nil)
)

type bit bool

const (
	zero bit = false
	one  bit = true
)

func newBStream(capacity int) *bstream {
	return newBStreamWithData(make([]byte, 0, capacity))
}

func newBStreamWithData(data []byte) *bstream {
	return &bstream{
		stream: data,
		wBit:   0,
		rIdx:   0,
		rBit:   8,
	}
}

type bstream struct {
	stream []byte
	wBit   uint8
	rIdx   int
	rBit   uint8
}

func (bs *bstream) clone() *bstream {
	stream := make([]byte, len(bs.stream))
	copy(stream, bs.stream)
	return &bstream{
		stream: stream,
		wBit:   bs.wBit,
		rIdx:   bs.rIdx,
		rBit:   bs.rBit,
	}
}

func (bs *bstream) rewind() {
	bs.rIdx = 0
	bs.rBit = 8
}

func (bs *bstream) writeBit(bit bit) {
	if bs.wBit == 0 {
		bs.stream = append(bs.stream, 0)
		bs.wBit = 8
	}

	if bit {
		bs.stream[len(bs.stream)-1] |= 1 << (bs.wBit - 1)
	}

	bs.wBit--
}

func (bs *bstream) writeByte(b byte) {
	if bs.wBit == 0 {
		bs.stream = append(bs.stream, b)
		return
	}

	// a byte must be append
	bs.stream = append(bs.stream, 0)
	bs.stream[len(bs.stream)-2] |= b >> (8 - bs.wBit)
	bs.stream[len(bs.stream)-1] |= b << bs.wBit
}

func (bs *bstream) writeBits(u uint64, nbits uint) {
	u <<= (64 - nbits)
	for nbits >= 8 {
		bs.writeByte(byte(u >> 56))
		u <<= 8
		nbits -= 8
	}

	for nbits > 0 {
		bs.writeBit((u >> 63) == 1)
		u <<= 1
		nbits--
	}
}

func (bs *bstream) checkReadBitSize() error {
	written := len(bs.stream)*8 - int(bs.wBit)
	read := (bs.rIdx+1)*8 - int(bs.rBit)
	if read >= written {
		return io.EOF
	}

	return nil
}

func (bs *bstream) readBit() (bit, error) {
	if err := bs.checkReadBitSize(); err != nil {
		return false, err
	}

	if bs.rBit == 0 {
		bs.rIdx++
		if bs.rIdx == len(bs.stream) {
			return false, io.EOF
		}

		bs.rBit = 8
	}

	bs.rBit--

	d := bs.stream[bs.rIdx] & (1 << bs.rBit)
	return d != 0, nil
}

func (bs *bstream) readByte() (byte, error) {
	if err := bs.checkReadBitSize(); err != nil {
		return 0, err
	}

	if bs.rBit == 0 {
		bs.rIdx++
		if bs.rIdx == len(bs.stream) {
			return 0, io.EOF
		}

		return bs.stream[bs.rIdx], nil
	}

	// we must read 8 bit data from 2 seperate bytes
	if (bs.rIdx + 1) == len(bs.stream) {
		return 0, io.EOF
	}

	bs.rIdx++
	return (bs.stream[bs.rIdx-1] << (8 - bs.rBit)) | (bs.stream[bs.rIdx] >> bs.rBit), nil
}

func (bs *bstream) readBits(nbits uint) (uint64, error) {
	var u uint64
	for nbits >= 8 {
		b, err := bs.readByte()
		if err != nil {
			return 0, err
		}

		u = (u << 8) | uint64(b)

		nbits -= 8
	}

	for nbits > 0 {
		bt, err := bs.readBit()
		if err != nil {
			return 0, err
		}
		u <<= 1
		if bt {
			u |= 1
		}

		nbits--
	}

	return u, nil
}

// MarshalBinary impl encoding.BinaryMarshaler
func (bs *bstream) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, bs.wBit); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, int64(bs.rIdx)); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, bs.rBit); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, bs.stream); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary impl encoding.BinaryUnmarshaler
func (bs *bstream) UnmarshalBinary(b []byte) error {
	buf := bytes.NewReader(b)

	if err := binary.Read(buf, binary.BigEndian, &bs.wBit); err != nil {
		return err
	}

	var rIdx int64
	if err := binary.Read(buf, binary.BigEndian, &rIdx); err != nil {
		return err
	}

	bs.rIdx = int(rIdx)

	if err := binary.Read(buf, binary.BigEndian, &bs.rBit); err != nil {
		return err
	}

	bs.stream = make([]byte, buf.Len())

	return binary.Read(buf, binary.BigEndian, &bs.stream)
}
