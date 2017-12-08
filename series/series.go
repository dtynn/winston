package series

import (
	"fmt"
	"sync"
	"time"

	"github.com/dgryski/go-bits"
)

const (
	defaultLeading = ^uint8(0)

	dodControlBits0    uint64 = 0x00
	dodControlBits10          = 0x02
	dodControlBits110         = 0x06
	dodControlBits1110        = 0x0e
	dodControlBits1111        = 0x0f

	valueControlBits0  uint64 = 0x00
	valueControlBits10        = 0x02
	valueControlBits11        = 0x03
)

// NewSeries return new series with second precision
func NewSeries(t time.Time) *Series {
	return newSeries(t, time.Second)
}

// NewMilliSeries return new series with millisecond precision
func NewMilliSeries(t time.Time) *Series {
	return newSeries(t, time.Millisecond)
}

func newSeries(t time.Time, precision time.Duration) *Series {
	s := &Series{
		t0:                t,
		bs:                newBStream(10240),
		leading:           defaultLeading,
		precision:         precision,
		precisionSettings: precisions[precision],
	}

	s.bs.writeBits(uint64(s.t0.UnixNano()), 64)
	return s
}

// Series time series
type Series struct {
	sync.RWMutex

	precision time.Duration
	precisionSettings

	t0        time.Time
	prevT     time.Time
	tdelta    int32
	prevVBits uint64

	bs       *bstream
	leading  uint8
	trailing uint8

	finished bool

	err error
}

func finish(bs *bstream, precision time.Duration) {
	// write an end-of-stream record
	bs.writeBits(dodControlBits1111, 4)
	bs.writeBits(precisions[precision].finish.bits, precisions[precision].finish.n)
	bs.writeBit(zero)
}

// Finish finish a stream
func (s *Series) Finish() {
	s.Lock()

	if !s.finished {
		finish(s.bs, s.precision)
		s.finished = true
	}

	s.Unlock()
}

// Push push timestamp and value bits
func (s *Series) Push(t time.Time, vbits uint64) {
	s.Lock()
	defer s.Unlock()

	t = t.Truncate(s.precision)

	if s.prevT.IsZero() {
		s.prevT = t
		s.prevVBits = vbits
		s.tdelta = int32(t.Sub(s.t0) / s.precision)
		// with one-day block, and precision of millisecond, we need at most 28 bits
		s.bs.writeBits(uint64(zigzag(s.tdelta)), 28)
		s.bs.writeBits(s.prevVBits, 64)
		return
	}

	// deal with delta-of-delta of timestamp
	tdelta := int32(t.Sub(s.prevT) / s.precision)
	// delta-of-delta
	dod := tdelta - s.tdelta
	var dodCtrlBits uint64

	// in the paper of facebook's gorilla,
	// they use
	//   '10' for [-63, 64]
	//  '110' for [-255, 256]
	// '1110' for [-2047, 2048]
	// '1111' for others
	switch {
	case dod == 0:
		// '0'
		s.bs.writeBit(zero)

	case inRange(dod, s.precisionSettings.dod[dodControlBits10].dodRange):

		dodCtrlBits = dodControlBits10
		s.bs.writeBits(dodControlBits10, 2)

	case inRange(dod, s.precisionSettings.dod[dodControlBits10].dodRange):

		dodCtrlBits = dodControlBits110
		s.bs.writeBits(dodControlBits110, 3)

	case inRange(dod, s.precisionSettings.dod[dodControlBits1110].dodRange):

		dodCtrlBits = dodControlBits1110
		s.bs.writeBits(dodControlBits1110, 4)

	default:
		// '1111' & 28 bit value
		dodCtrlBits = dodControlBits1111
		s.bs.writeBits(dodControlBits1111, 4)
	}

	if dodCtrlBits > 0 {
		s.bs.writeBits(uint64(zigzag(dod)), s.precisionSettings.dod[dodCtrlBits].dodNBits)
	}

	vdelta := vbits ^ s.prevVBits
	if vdelta == 0 {
		s.bs.writeBit(zero)
	} else {
		s.bs.writeBit(one)

		// When XOR is non-zero, calculate the number of leading
		// and trailing zeros in the XOR, store bit ‘1’ followed
		// by either a) or b):
		// 		(a) (Control bit ‘0’) If the block of meaningful bits
		//           falls within the block of previous meaningful bits,
		//           i.e., there are at least as many leading zeros and
		//           as many trailing zeros as with the previous value,
		// 		     use that information for the block position and
		//           just store the meaningful XORed value.
		// 		(b) (Control bit ‘1’) Store the length of the number
		//           of leading zeros in the next 5 bits, then store the
		//           length of the meaningful XORed value in the next
		//           6 bits. Finally store the meaningful bits of the
		//           XORed value.

		leading := uint8(bits.Clz(vdelta))
		trailing := uint8(bits.Ctz(vdelta))

		// leading has been set and for the meaningful bit, previous size >= current size
		if s.leading != defaultLeading && leading >= s.leading && trailing >= s.trailing {
			s.bs.writeBit(zero)
			s.bs.writeBits(vdelta>>s.trailing, uint(64-s.leading-s.trailing))
		} else {
			s.leading, s.trailing = leading, trailing

			s.bs.writeBit(one)
			// we use 6 bit for storing leading size to support at most 63 bit leading zeros
			s.bs.writeBits(uint64(leading), 6)

			meaningfulBits := 64 - leading - trailing
			s.bs.writeBits(uint64(meaningfulBits), 6)
			s.bs.writeBits(vdelta>>trailing, uint(meaningfulBits))
		}
	}

	s.prevT = t
	s.tdelta = tdelta
	s.prevVBits = vbits
}

// Iter return an iterator
func (s *Series) Iter() (*Iter, error) {
	s.Lock()
	bs := s.bs.clone()
	s.Unlock()

	finish(bs, s.precision)
	bs.rewind()

	return bstreamIter(bs, s.precision)
}

func bstreamIter(bs *bstream, precision time.Duration) (*Iter, error) {
	precSettings, ok := precisions[precision]
	if !ok {
		return nil, fmt.Errorf("unsupport precision %s", precision)
	}

	t0bits, err := bs.readBits(64)
	if err != nil {
		return nil, err
	}

	it := &Iter{
		t0:                time.Unix(0, int64(t0bits)),
		bs:                bs,
		precision:         precision,
		precisionSettings: precSettings,
	}

	it.Stat.dod = map[uint64]int{}
	it.Stat.vdelta = map[uint64]int{}

	return it, nil
}

// NewIter return new iterator with given data
func NewIter(data []byte, precision time.Duration) (*Iter, error) {
	bs := newBStreamWithData(data)
	return bstreamIter(bs, precision)
}

// Iter stream iter
type Iter struct {
	precision time.Duration
	precisionSettings

	t0     time.Time
	t      time.Time
	tdelta int32
	vbits  uint64

	bs       *bstream
	leading  uint8
	trailing uint8

	finished bool

	pointStat bool

	Stat struct {
		points int
		dod    map[uint64]int
		vdelta map[uint64]int
	}

	err error
}

// PointStat if we stat the points
func (i *Iter) PointStat(b bool) {
	i.pointStat = b
}

// Next try read next value
func (i *Iter) Next() bool {
	if i.err != nil || i.finished {
		return false
	}

	if i.t.IsZero() {
		tdeltabits, err := i.bs.readBits(28)
		if err != nil {
			i.err = fmt.Errorf("read first tdelta: %s", err)
			return false
		}

		vbits, err := i.bs.readBits(64)
		if err != nil {
			i.err = fmt.Errorf("read first value bits: %s", err)
			return false
		}

		i.tdelta = zagzig(uint32(tdeltabits))
		i.t = newTime(i.t0, i.tdelta, i.precision)
		i.vbits = vbits

		if i.pointStat {
			i.Stat.points++
		}

		return true
	}

	dodCtrlBits, err := readDoDControlBits(i.bs)
	if err != nil {
		i.err = fmt.Errorf("read dod control bits: %s", err)
		return false
	}

	var dodNBits uint

	switch dodCtrlBits {
	case dodControlBits0:
		// dodNBits = 0
		// i.tdelta = i.tdelta

	case dodControlBits10,
		dodControlBits110,
		dodControlBits1110,
		dodControlBits1111:
		dodNBits = i.precisionSettings.dod[dodCtrlBits].dodNBits

	default:
		i.err = fmt.Errorf("malformed dod control bits %02x", dodCtrlBits)
		return false

	}

	var dod int32
	if dodNBits > 0 {
		dodbits, err := i.bs.readBits(dodNBits)
		if err != nil {
			i.err = fmt.Errorf("read dod bits: %s", err)
			return false
		}

		if dodCtrlBits == dodControlBits1111 && dodbits == i.precisionSettings.finish.bits {
			i.finished = true
			return false
		}

		dod = zagzig(uint32(dodbits))
	}

	i.tdelta += dod
	i.t = newTime(i.t, i.tdelta, i.precision)

	valCtrlBits, err := readValueControlBits(i.bs)
	if err != nil {
		i.err = fmt.Errorf("read value control bits: %s", err)
		return false
	}

	var vdelta uint64

	switch valCtrlBits {
	case valueControlBits0:
		// vdelta = 0

	case valueControlBits10:
		meaningfulNBits := uint(64 - i.leading - i.trailing)
		meaningful, err := i.bs.readBits(meaningfulNBits)
		if err != nil {
			i.err = fmt.Errorf("read meaningful value with control bits %02x: %s", valCtrlBits, err)
			return false
		}

		vdelta = meaningful << i.trailing

	case valueControlBits11:
		leading, err := i.bs.readBits(6)
		if err != nil {
			i.err = fmt.Errorf("read leading bits: %s", err)
			return false
		}

		meaningfulNbits, err := i.bs.readBits(6)
		if err != nil {
			i.err = fmt.Errorf("read meaningful bits: %s", err)
			return false
		}

		meaningful, err := i.bs.readBits(uint(meaningfulNbits))
		if err != nil {
			i.err = fmt.Errorf("read meaningful value with control bits %02x: %s", valCtrlBits, err)
			return false
		}

		i.leading = uint8(leading)
		i.trailing = uint8(64 - leading - meaningfulNbits)
		vdelta = meaningful << i.trailing

	default:
		i.err = fmt.Errorf("malformed value control bits %02x", valCtrlBits)
		return false
	}

	i.vbits ^= vdelta

	if i.pointStat {
		i.Stat.points++
		i.Stat.dod[dodCtrlBits]++
		i.Stat.vdelta[valCtrlBits]++
	}

	return true
}

// Err return last error
func (i *Iter) Err() error {
	return i.err
}

// Point return current point
func (i *Iter) Point() (time.Time, uint64) {
	return i.t, i.vbits
}

func newTime(old time.Time, tdelta int32, precision time.Duration) time.Time {
	return old.Add(time.Duration(tdelta) * precision)
}

func readDoDControlBits(bs *bstream) (uint64, error) {
	var bits uint64
	for read := 0; read < 4; read++ {
		b, err := bs.readBit()
		if err != nil {
			return 0, err
		}

		bits <<= 1
		// get zero bit, break
		if !b {
			break
		}

		bits |= 1
	}

	return bits, nil
}

func readValueControlBits(bs *bstream) (uint64, error) {
	var bits uint64
	for read := 0; read < 2; read++ {
		b, err := bs.readBit()
		if err != nil {
			return 0, err
		}

		bits <<= 1
		// get zero bit, break
		if !b {
			break
		}

		bits |= 1
	}

	return bits, nil
}

func inRange(val int32, r [2]int32) bool {
	return r[0] <= val && val <= r[1]
}