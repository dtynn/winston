package zigzag

import (
	"math"
	"math/rand"
	"testing"
)

func TestZigZag(t *testing.T) {
	t.Run("zigzag32", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			positive := rand.Int31n(math.MaxInt32)
			negative := -positive

			if got := ZagZig32(ZigZag32(positive)); got != positive {
				t.Errorf("expected positive value %d, got %d", positive, got)
			}

			if got := ZagZig32(ZigZag32(negative)); got != negative {
				t.Errorf("expected negative value %d, got %d", negative, got)
			}
		}
	})

	t.Run("zigzag", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			positive := rand.Int63n(math.MaxInt64)
			negative := -positive

			if got := ZagZig(ZigZag(positive)); got != positive {
				t.Errorf("expected positive value %d, got %d", positive, got)
			}

			if got := ZagZig(ZigZag(negative)); got != negative {
				t.Errorf("expected negative value %d, got %d", negative, got)
			}
		}
	})
}
