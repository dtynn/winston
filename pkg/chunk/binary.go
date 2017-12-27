package chunk

import (
	"encoding/binary"
	"io"
)

type bwriter struct {
	io.Writer
	err error
}

func (b *bwriter) write(data interface{}) {
	if b.err == nil {
		b.err = binary.Write(b.Writer, binary.BigEndian, data)
	}
}

type breader struct {
	io.Reader
	err error
}

func (b *breader) read(data interface{}) {
	if b.err == nil {
		b.err = binary.Read(b.Reader, binary.BigEndian, data)
	}
}
