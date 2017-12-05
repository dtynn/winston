package winston

import (
	"io"
	"math/rand"
	"testing"
	"time"
)

func TestBStreamWriteAndReadBit(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := 7 + rand.Intn(13)

	sets := make([]struct {
		bit bit
	}, n)

	for i := 0; i < n; i++ {
		sets[i].bit = rand.Intn(2) == 1
	}

	t.Logf("n is %d", n)
	bs := newBStream(10)

	for _, s := range sets {
		bs.writeBit(s.bit)
	}

	t.Logf("%#v", bs.stream)

	for i, s := range sets {
		b, err := bs.readBit()
		if err != nil {
			t.Fatalf("#%d expected nil error, got %s", i+1, err)
		}

		if b != s.bit {
			t.Errorf("#%d expected %v, got %v", i+1, s.bit, b)
		}
	}

	_, err := bs.readBit()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v: %d %d %d %d", err, len(bs.stream), bs.wBit, bs.rIdx, bs.rBit)
	}
}

func TestBStreamWriteAndReadByte(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	n := 7 + rand.Intn(13)

	sets := make([]struct {
		b byte
	}, n)

	for i := 0; i < n; i++ {
		sets[i].b = byte(rand.Intn(256))
	}

	t.Logf("n is %d", n)
	bs := newBStream(50)

	for i := range sets {
		if i%3 == 0 {
			bs.writeBit(i%2 == 0)
		}

		bs.writeByte(sets[i].b)
	}

	t.Logf("%#v", bs.stream)

	for i := range sets {
		if i%3 == 0 {
			bit, err := bs.readBit()
			if err != nil {
				t.Fatalf("#%d expected nil error reading bit, got %s", i+1, err)
			}

			if bit != (i%2 == 0) {
				t.Errorf("#%d expected bit %v, got %v", i+1, (i%2 == 0), bit)
			}
		}

		b, err := bs.readByte()
		if err != nil {
			t.Fatalf("#%d expected nil error reading byte, got %s", i+1, err)
		}

		if b != sets[i].b {
			t.Errorf("#%d expcted byte %v, got %v", i+1, sets[i].b, b)
		}
	}

	_, err := bs.readBit()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}

func TestBStreamWriteAndReadBits(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	n := 7 + rand.Intn(13)
	sets := make([]struct {
		u     uint64
		nbits uint
	}, n)

	now := uint64(time.Now().UnixNano())

	for i := 0; i < n; i++ {
		sets[i].u = now + uint64(rand.Int63n(int64(time.Millisecond)))
		sets[i].nbits = 64 - uint(rand.Intn(3))
	}

	t.Logf("n is %d", n)

	bs := newBStream(50)

	for _, s := range sets {
		bs.writeBits(s.u, s.nbits)
	}

	t.Logf("%#v", bs.stream)

	for i, s := range sets {
		read, err := bs.readBits(s.nbits)
		if err != nil {
			t.Fatalf("#%d expected nil error, got %s, pos: %d, %d", i+1, err, bs.rIdx, bs.rBit)
		}

		if read != s.u {
			t.Errorf("#%d expected %d, got %d", i+1, s.u, read)
		}
	}

	_, err := bs.readBit()
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}

	// test MarshalBinary
	data, err := bs.MarshalBinary()
	if err != nil {
		t.Fatalf("expeted nil error marshaling binary, got %s", err)
	}

	newBS := new(bstream)
	if err := newBS.UnmarshalBinary(data); err != nil {
		t.Fatalf("expected nil error unmarshaling binary, got %s", err)
	}

	if newBS.wBit != bs.wBit {
		t.Fatalf("malformed wBit %d - %d", bs.wBit, newBS.wBit)
	}

	// if newBS.rIdx != bs.rIdx {
	// 	t.Fatalf("malformed rIdx %d - %d", bs.rIdx, newBS.rIdx)
	// }

	// if newBS.rBit != bs.rBit {
	// 	t.Fatalf("malformed rBit %d - %d", bs.rBit, newBS.rBit)
	// }

	if len(newBS.stream) != len(bs.stream) {
		t.Fatalf("malformed stream length %d - %d", len(newBS.stream), len(bs.stream))
	}

	for i := range newBS.stream {
		if newBS.stream[i] != bs.stream[i] {

		}
	}
}
