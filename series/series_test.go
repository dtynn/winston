package series

import (
	"math/rand"
	"testing"
	"time"
)

func TestReadDodControlBits(t *testing.T) {
	bits := []struct {
		ctrlBits uint64
		nbits    uint
	}{
		{
			ctrlBits: dodControlBits10,
			nbits:    2,
		},
		{
			ctrlBits: dodControlBits1110,
			nbits:    4,
		},
		{
			ctrlBits: dodControlBits110,
			nbits:    3,
		},
		{
			ctrlBits: dodControlBits1111,
			nbits:    4,
		},
		{
			ctrlBits: dodControlBits0,
			nbits:    1,
		},
	}

	bs := newBStream(50)

	for _, b := range bits {
		bs.writeBits(b.ctrlBits, b.nbits)
	}

	for i, b := range bits {
		ctrlBits, err := readDoDControlBits(bs)
		if err != nil {
			t.Fatalf("#%d got read err: %s", i+1, err)
		}

		if ctrlBits != b.ctrlBits {
			t.Errorf("#%d expcted ctrl bits %02x, got %02x", i+1, b.ctrlBits, ctrlBits)
		}
	}
}

func TestMilliSeriesAndIter(t *testing.T) {
	baseT := time.Now()
	ti := baseT.Add(time.Hour)

	points := make([]struct {
		t   time.Time
		val uint64
	}, 2000000)

	for i := range points {
		n := uint(4 + i%20)
		ti = ti.Add(time.Duration(rand.Int31n(1<<n)) * time.Millisecond)
		points[i].t = ti
		points[i].val = uint64(6 + rand.Int63n(14))
	}

	ts := NewMilliSeries(baseT.Truncate(24 * time.Hour))

	for i := range points {
		ts.Push(points[i].t, points[i].val)
	}

	iter, err := ts.Iter()
	if err != nil {
		t.Fatalf("new iterator: %s", err)
	}

	iter.PointStat(true)

	i := 0
	for {
		if !iter.Next() {
			break
		}

		pt, pv := iter.Point()
		expectedT := points[i].t.Truncate(time.Millisecond)
		if pt != expectedT {
			t.Fatalf("#%d expected point time %s, got %s", i+1, expectedT, pt)
		}

		if pv != points[i].val {
			t.Fatalf("#%d expected point val %d, got %d", i+1, points[i].val, pv)
		}

		i++
	}

	if i != len(points) {
		t.Fatalf("expected %d points, got %d", len(points), i)
	}

	t.Logf("expected %d Bytes, got %d Bytes", len(points)*16, len(ts.bs.stream))
	t.Logf("point stat %v", iter.Stat)
}

func TestSeriesAndIter(t *testing.T) {
	baseT := time.Now()
	ti := baseT.Add(time.Hour)

	points := make([]struct {
		t   time.Time
		val uint64
	}, 1000000)

	for i := range points {
		n := uint(4 + i%10)
		ti = ti.Add(time.Duration(rand.Int31n(1<<n)) * time.Second)
		points[i].t = ti
		points[i].val = uint64(6 + rand.Int63n(14))
	}

	ts := NewSeries(baseT.Truncate(24 * time.Hour))

	for i := range points {
		ts.Push(points[i].t, points[i].val)
	}

	iter, err := ts.Iter()
	if err != nil {
		t.Fatalf("new iterator: %s", err)
	}

	iter.PointStat(true)

	i := 0
	for {
		if !iter.Next() {
			break
		}

		pt, pv := iter.Point()
		expectedT := points[i].t.Truncate(time.Second)
		if pt != expectedT {
			t.Fatalf("#%d expected point time %s, got %s", i+1, expectedT, pt)
		}

		if pv != points[i].val {
			t.Fatalf("#%d expected point val %d, got %d", i+1, points[i].val, pv)
		}

		i++
	}

	if i != len(points) {
		t.Fatalf("expected %d points, got %d", len(points), i)
	}

	t.Logf("expected %d Bytes, got %d Bytes", len(points)*16, len(ts.bs.stream))
	t.Logf("point stat %v", iter.Stat)
}

func BenchmarkSeriesPush(b *testing.B) {
	baseT := time.Now()
	tm := baseT.Add(time.Hour)

	ts := NewMilliSeries(baseT.Truncate(24 * time.Hour))

	for i := 0; i < b.N; i++ {
		tm = tm.Add(time.Duration(350+rand.Int63n(300)) * time.Millisecond)
		ts.Push(tm, uint64(6+rand.Int63n(14)))
	}
}
