package series

import (
	"fmt"
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

func TestZigZagBitSize(t *testing.T) {
	for n := 0; n < 31; n++ {
		var i int32 = 1 << uint(n)
		posstr := fmt.Sprintf("%b", zigzag(i))
		negstr := fmt.Sprintf("%b", zigzag(-i))

		if len(posstr) != n+2 {
			t.Errorf("for (1 << %d), expected %d bit, got %d", n, n+2, len(posstr))
		}

		if len(negstr) != n+1 {
			t.Errorf("for -(1 << %d), expected %d bit, got %d", n, n+1, len(negstr))
		}
	}
}
