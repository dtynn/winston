package chunk

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

func TestMilliChunkAndIter(t *testing.T) {
	var pointNum uint64 = 1000000

	baseT := time.Now()
	ti := baseT.Add(time.Hour)

	points := make([]struct {
		t   time.Time
		val uint64
	}, pointNum)

	for i := range points {
		n := uint(4 + i%20)
		ti = ti.Add(time.Duration(rand.Int31n(1<<n)) * time.Millisecond)
		points[i].t = ti
		points[i].val = uint64(6 + rand.Int63n(14))
	}

	ts := NewMilliChunk(baseT.Truncate(24 * time.Hour))

	for i := range points {
		ts.PushTime(points[i].t, points[i].val)
	}

	iter, err := ts.Iter()
	if err != nil {
		t.Fatalf("new iterator: %s", err)
	}

	if iter.num != pointNum {
		t.Fatalf("expected total points num %d, got %d", pointNum, iter.num)
	}

	iter.PointStat(true)

	i := 0
	for {
		if !iter.Next() {
			break
		}

		pt, pv := iter.Point()
		expectedT := points[i].t.Truncate(time.Millisecond)
		gotT := iter.PointTime(pt)
		if gotT != expectedT {
			t.Fatalf("#%d expected point time %s, got %s", i+1, expectedT, gotT)
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

func TestMilliChunkAndIterWithMarshaling(t *testing.T) {
	var pointNum uint64 = 1000000
	half := pointNum / 2

	baseT := time.Now()
	ti := baseT.Add(time.Hour)

	points := make([]struct {
		t   time.Time
		val uint64
	}, pointNum)

	for i := range points {
		n := uint(4 + i%20)
		ti = ti.Add(time.Duration(rand.Int31n(1<<n)) * time.Millisecond)
		points[i].t = ti
		points[i].val = uint64(6 + rand.Int63n(14))
	}

	ts := NewMilliChunk(baseT.Truncate(24 * time.Hour))

	for i := range points {
		ts.PushTime(points[i].t, points[i].val)

		if uint64(i) == half {
			data, err := ts.MarshalBinary()
			if err != nil {
				t.Fatalf("marshal binary %s", err)
			}

			ck := new(Chunk)
			if err := ck.UnmarshalBinary(data); err != nil {
				t.Fatalf("unmarshal binary %s", err)
			}

			ts = ck
		}
	}

	iter, err := ts.Iter()
	if err != nil {
		t.Fatalf("new iterator: %s", err)
	}

	if iter.num != pointNum {
		t.Fatalf("expected total points num %d, got %d", pointNum, iter.num)
	}

	iter.PointStat(true)

	i := 0
	for {
		if !iter.Next() {
			break
		}

		pt, pv := iter.Point()
		expectedT := points[i].t.Truncate(time.Millisecond)
		gotT := iter.PointTime(pt)
		if gotT != expectedT {
			t.Fatalf("#%d expected point time %s, got %s", i+1, expectedT, gotT)
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

func TestChunkAndIter(t *testing.T) {
	var pointNum uint64 = 1000000

	baseT := time.Now()
	ti := baseT.Add(time.Hour)

	points := make([]struct {
		t   time.Time
		val uint64
	}, pointNum)

	for i := range points {
		n := uint(4 + i%10)
		ti = ti.Add(time.Duration(rand.Int31n(1<<n)) * time.Second)
		points[i].t = ti
		points[i].val = uint64(6 + rand.Int63n(14))
	}

	ts := NewChunk(baseT.Truncate(24 * time.Hour))

	for i := range points {
		ts.PushTime(points[i].t, points[i].val)
	}

	iter, err := ts.Iter()
	if err != nil {
		t.Fatalf("new iterator: %s", err)
	}

	if iter.num != pointNum {
		t.Fatalf("expected total points num %d, got %d", pointNum, iter.num)
	}

	iter.PointStat(true)

	i := 0
	for {
		if !iter.Next() {
			break
		}

		pt, pv := iter.Point()
		expectedT := points[i].t.Truncate(time.Second)
		gotT := iter.PointTime(pt)
		if gotT != expectedT {
			t.Fatalf("#%d expected point time %s, got %s", i+1, expectedT, gotT)
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

func BenchmarkChunkPush(b *testing.B) {
	baseT := time.Now()
	tm := baseT.Add(time.Hour)

	ts := NewMilliChunk(baseT.Truncate(24 * time.Hour))

	for i := 0; i < b.N; i++ {
		tm = tm.Add(time.Duration(350+rand.Int63n(300)) * time.Millisecond)
		ts.PushTime(tm, uint64(6+rand.Int63n(14)))
	}
}

func BenchmarkChunkIterRead(b *testing.B) {
	b.StopTimer()
	baseT := time.Now()
	tm := baseT.Add(time.Hour)

	ts := NewMilliChunk(baseT.Truncate(24 * time.Hour))

	for i := 0; i < b.N; i++ {
		tm = tm.Add(time.Duration(350+rand.Int63n(300)) * time.Millisecond)
		ts.PushTime(tm, uint64(6+rand.Int63n(14)))
	}

	iter, err := ts.Iter()
	if err != nil {
		b.Fatalf("get iter %s", err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if !iter.Next() {
			b.Fatalf("read through all points at benchmark loop %d", b.N)
		}
	}
}
