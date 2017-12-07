package winston

import (
	"math"
	"time"
)

type precisionSettings struct {
	dod map[uint64]struct {
		dodRange [2]int32
		dodNBits uint
	}

	finish struct {
		bits uint64
		n    uint
	}
}

var precisions = map[time.Duration]precisionSettings{
	time.Millisecond: {
		dod: map[uint64]struct {
			dodRange [2]int32
			dodNBits uint
		}{
			dodControlBits10: {
				[2]int32{-512, 511},
				10,
			},
			dodControlBits110: {
				[2]int32{-32768, 32767},
				16,
			},
			dodControlBits1110: {
				[2]int32{-2097152, 2097151},
				22,
			},
			dodControlBits1111: {
				[2]int32{math.MinInt32, math.MaxInt32},
				28,
			},
		},

		finish: struct {
			bits uint64
			n    uint
		}{
			^uint64(0) >> 36,
			28,
		},
	},

	time.Second: {
		dod: map[uint64]struct {
			dodRange [2]int32
			dodNBits uint
		}{
			dodControlBits10: {
				[2]int32{-64, 63},
				7,
			},
			dodControlBits110: {
				[2]int32{-256, 255},
				9,
			},
			dodControlBits1110: {
				[2]int32{-2048, 2047},
				12,
			},
			dodControlBits1111: {
				[2]int32{math.MinInt32, math.MaxInt32},
				32,
			},
		},

		finish: struct {
			bits uint64
			n    uint
		}{
			^uint64(0) >> 32,
			32,
		},
	},
}
