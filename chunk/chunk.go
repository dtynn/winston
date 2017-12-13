package chunk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/bits"
	"sync"
	"time"
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

// NewChunk return new series with second precision
func NewChunk(t time.Time) *Chunk {
	return newChunk(t, time.Second)
}

// NewMilliChunk return new series with millisecond precision
func NewMilliChunk(t time.Time) *Chunk {
	return newChunk(t, time.Millisecond)
}

func newChunk(t time.Time, precision time.Duration) *Chunk {
	c := &Chunk{
		t0:                t.UnixNano() / int64(precision),
		bs:                newBStream(10240),
		leading:           defaultLeading,
		precision:         precision,
		precisionSettings: precisions[precision],
	}

	// placeholder for point num
	c.bs.writeBits(0, 64)
	c.bs.writeBits(uint64(c.t0), 64)
	return c
}

// Chunk time series
type Chunk struct {
	sync.RWMutex

	precision time.Duration
	precisionSettings

	t0        int64
	prevT     int64
	tdelta    int32
	prevVBits uint64

	bs       *bstream
	leading  uint8
	trailing uint8

	finished bool

	num uint64
}

func finish(bs *bstream, precision time.Duration) {
	// write an end-of-stream record
	bs.writeBits(dodControlBits1111, 4)
	bs.writeBits(precisions[precision].finish.bits, precisions[precision].finish.n)
	bs.writeBit(zero)
}

// Finish finish a stream
func (c *Chunk) Finish() {
	c.Lock()

	if !c.finished {
		finish(c.bs, c.precision)
		c.finished = true
	}

	c.Unlock()
}

// PushTime push time.Time and value bits
func (c *Chunk) PushTime(t time.Time, vbits uint64) {
	c.Push(t.UnixNano()/int64(c.precision), vbits)
}

// Push push timestamp and value bits
func (c *Chunk) Push(t int64, vbits uint64) {
	c.Lock()
	defer c.Unlock()

	c.num++
	binary.BigEndian.PutUint64(c.bs.stream[:8], c.num)

	if c.prevT == 0 {
		c.prevT = t
		c.prevVBits = vbits
		c.tdelta = int32(t - c.t0)

		// with one-day block, and precision of millisecond, we need at most 28 bits
		c.bs.writeBits(uint64(zigzag(c.tdelta)), 28)
		c.bs.writeBits(c.prevVBits, 64)
		return
	}

	// deal with delta-of-delta of timestamp
	tdelta := int32(t - c.prevT)
	// delta-of-delta
	dod := tdelta - c.tdelta
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
		c.bs.writeBit(zero)

	case inRange(dod, c.precisionSettings.dod[dodControlBits10].dodRange):

		dodCtrlBits = dodControlBits10
		c.bs.writeBits(dodControlBits10, 2)

	case inRange(dod, c.precisionSettings.dod[dodControlBits10].dodRange):

		dodCtrlBits = dodControlBits110
		c.bs.writeBits(dodControlBits110, 3)

	case inRange(dod, c.precisionSettings.dod[dodControlBits1110].dodRange):

		dodCtrlBits = dodControlBits1110
		c.bs.writeBits(dodControlBits1110, 4)

	default:
		// '1111' & 28 bit value
		dodCtrlBits = dodControlBits1111
		c.bs.writeBits(dodControlBits1111, 4)
	}

	if dodCtrlBits > 0 {
		c.bs.writeBits(uint64(zigzag(dod)), c.precisionSettings.dod[dodCtrlBits].dodNBits)
	}

	vdelta := vbits ^ c.prevVBits
	if vdelta == 0 {
		c.bs.writeBit(zero)
	} else {
		c.bs.writeBit(one)

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

		leading := uint8(bits.LeadingZeros64(vdelta))
		trailing := uint8(bits.TrailingZeros64(vdelta))

		// leading has been set and for the meaningful bit, previous size >= current size
		if c.leading != defaultLeading && leading >= c.leading && trailing >= c.trailing {
			c.bs.writeBit(zero)
			c.bs.writeBits(vdelta>>c.trailing, uint(64-c.leading-c.trailing))
		} else {
			c.leading, c.trailing = leading, trailing

			c.bs.writeBit(one)
			// we use 6 bit for storing leading size to support at most 63 bit leading zeros
			c.bs.writeBits(uint64(leading), 6)

			meaningfulBits := 64 - leading - trailing
			c.bs.writeBits(uint64(meaningfulBits), 6)
			c.bs.writeBits(vdelta>>trailing, uint(meaningfulBits))
		}
	}

	c.prevT = t
	c.tdelta = tdelta
	c.prevVBits = vbits
}

// MarshalBinary impl encoding.BinaryMarshaler
func (c *Chunk) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := bwrite(buf, c.precision); err != nil {
		return nil, err
	}

	if err := bwrite(buf, c.t0); err != nil {
		return nil, err
	}

	if err := bwrite(buf, c.prevT); err != nil {
		return nil, err
	}

	if err := bwrite(buf, c.tdelta); err != nil {
		return nil, err
	}

	if err := bwrite(buf, c.prevVBits); err != nil {
		return nil, err
	}

	if err := bwrite(buf, c.leading); err != nil {
		return nil, err
	}

	if err := bwrite(buf, c.leading); err != nil {
		return nil, err
	}

	if err := bwrite(buf, c.finished); err != nil {
		return nil, err
	}

	if err := bwrite(buf, c.num); err != nil {
		return nil, err
	}

	bsdata, err := c.bs.MarshalBinary()
	if err != nil {
		return nil, err
	}

	if err := bwrite(buf, bsdata); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary impl encoding.BinaryUnmarshaler
func (c *Chunk) UnmarshalBinary(buf []byte) error {
	r := bytes.NewBuffer(buf)

	var prec int64

	if err := bread(r, &prec); err != nil {
		return err
	}

	c.precision = time.Duration(prec)
	settings, ok := precisions[c.precision]
	if !ok {
		return fmt.Errorf("malformed precision %s", c.precision)
	}
	c.precisionSettings = settings

	if err := bread(r, &c.t0); err != nil {
		return err
	}

	if err := bread(r, &c.prevT); err != nil {
		return err
	}

	if err := bread(r, &c.tdelta); err != nil {
		return err
	}

	if err := bread(r, &c.prevVBits); err != nil {
		return err
	}

	if err := bread(r, &c.leading); err != nil {
		return err
	}

	if err := bread(r, &c.leading); err != nil {
		return err
	}

	if err := bread(r, &c.finished); err != nil {
		return err
	}

	if err := bread(r, &c.num); err != nil {
		return err
	}

	bsdata := make([]byte, r.Len())
	if err := bread(r, &bsdata); err != nil {
		return err
	}

	c.bs = new(bstream)

	return c.bs.UnmarshalBinary(bsdata)
}

// Iter return an iterator
func (c *Chunk) Iter() (*Iter, error) {
	c.Lock()
	bs := c.bs.clone()
	c.Unlock()

	finish(bs, c.precision)
	bs.rewind()

	return bstreamIter(bs, c.precision)
}

func bstreamIter(bs *bstream, precision time.Duration) (*Iter, error) {
	precSettings, ok := precisions[precision]
	if !ok {
		return nil, fmt.Errorf("unsupport precision %s", precision)
	}

	_, err := bs.readBits(64)
	if err != nil {
		return nil, err
	}

	t0bits, err := bs.readBits(64)
	if err != nil {
		return nil, err
	}

	it := &Iter{
		t0:                int64(t0bits),
		bs:                bs,
		precision:         precision,
		precisionSettings: precSettings,
		num:               binary.BigEndian.Uint64(bs.stream[:8]),
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

	t0     int64
	t      int64
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

	num uint64

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

	if i.t == 0 {
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
		i.t = i.t0 + int64(i.tdelta)
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
	i.t = i.t + int64(i.tdelta)

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

// Total total points
func (i *Iter) Total() uint64 {
	return i.num
}

// Err return last error
func (i *Iter) Err() error {
	return i.err
}

// Point return current point
func (i *Iter) Point() (int64, uint64) {
	return i.t, i.vbits
}

// PointTime return point time from timestamp
func (i *Iter) PointTime(ts int64) time.Time {
	return time.Unix(0, ts*int64(i.precision))
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
