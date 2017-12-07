package series

import (
	"math"
	"math/rand"
	"testing"
)

func TestZigZag(t *testing.T) {
	for i := 0; i < 20; i++ {
		pos := rand.Int31n(math.MaxInt32)
		neg := -pos

		posv := zagzig(zigzag(pos))
		negv := zagzig(zigzag(neg))

		if pos != posv {
			t.Errorf("#%d expected positive value %d, got %d", i+1, pos, posv)
		}

		if neg != negv {
			t.Errorf("#%d expected negative value %d, got %d", i+1, neg, negv)
		}
	}
}
