package chunk

import (
	"encoding/binary"
	"io"
)

func bwrite(buf io.Writer, data interface{}) error {
	return binary.Write(buf, binary.BigEndian, data)
}

func bread(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.BigEndian, data)
}
