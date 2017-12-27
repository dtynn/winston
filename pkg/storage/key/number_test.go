package key

import (
	"math"
	"math/rand"
	"testing"
)

func TestNumber(t *testing.T) {
	prefix := []byte("_prefix_")

	t.Run("UI64", func(t *testing.T) {
		start := rand.Uint64()
		for i := 0; i < 20; i++ {
			n := UI64(start) + UI64(i)
			key := Key(prefix, n)
			var got UI64
			if err := Unmarshal(key, prefix, &got); err != nil {
				t.Fatalf("unexpected unmarshal error for %d: %s", n, err)
			}

			if got != n {
				t.Fatalf("expected %d, got %d", n, got)
			}

			if err := Unmarshal(key[:len(key)-2], prefix, &got); err != ErrMalformedKeySize {
				t.Fatalf("expected key size error, got %s", err)
			}

			if err := Unmarshal(append(key, 'a'), prefix, &got); err != ErrMalformedKeySize {
				t.Fatalf("expected key size error, got %s", err)
			}
		}
	})

	t.Run("I64", func(t *testing.T) {
		start := rand.Int63()
		for i := 0; i < 20; i++ {
			n := I64(start) + I64(i)
			key := Key(prefix, n)
			var got I64
			if err := Unmarshal(key, prefix, &got); err != nil {
				t.Fatalf("unexpected unmarshal error for %d: %s", n, err)
			}

			if got != n {
				t.Fatalf("expected %d, got %d", n, got)
			}

			if err := Unmarshal(key[:len(key)-2], prefix, &got); err != ErrMalformedKeySize {
				t.Fatalf("expected key size error, got %s", err)
			}

			if err := Unmarshal(append(key, 'a'), prefix, &got); err != ErrMalformedKeySize {
				t.Fatalf("expected key size error, got %s", err)
			}
		}
	})

	t.Run("UI32", func(t *testing.T) {
		start := rand.Uint32()
		for i := 0; i < 20; i++ {
			n := UI32(start) + UI32(i)
			key := Key(prefix, n)
			var got UI32
			if err := Unmarshal(key, prefix, &got); err != nil {
				t.Fatalf("unexpected unmarshal error for %d: %s", n, err)
			}

			if got != n {
				t.Fatalf("expected %d, got %d", n, got)
			}

			if err := Unmarshal(key[:len(key)-2], prefix, &got); err != ErrMalformedKeySize {
				t.Fatalf("expected key size error, got %s", err)
			}

			if err := Unmarshal(append(key, 'a'), prefix, &got); err != ErrMalformedKeySize {
				t.Fatalf("expected key size error, got %s", err)
			}

			// size > MaxVarintLen32
			overflowKey := Key(prefix, UI64(math.MaxUint64))
			if err := Unmarshal(overflowKey, prefix, &got); err != ErrMalformedKeySize {
				t.Fatalf("expected key size error for overflow number, got %s", err)
			}
		}
	})

	t.Run("I32", func(t *testing.T) {
		start := rand.Int31()
		for i := 0; i < 20; i++ {
			n := I32(start) + I32(i)
			key := Key(prefix, n)
			var got I32
			if err := Unmarshal(key, prefix, &got); err != nil {
				t.Fatalf("unexpected unmarshal error for %d: %s", n, err)
			}

			if got != n {
				t.Fatalf("expected %d, got %d", n, got)
			}

			if err := Unmarshal(key[:len(key)-2], prefix, &got); err != ErrMalformedKeySize {
				t.Fatalf("expected key size error, got %s", err)
			}

			if err := Unmarshal(append(key, 'a'), prefix, &got); err != ErrMalformedKeySize {
				t.Fatalf("expected key size error, got %s", err)
			}

			// size > MaxVarintLen32
			overflowKey := Key(prefix, I64(math.MaxInt64))
			if err := Unmarshal(overflowKey, prefix, &got); err != ErrMalformedKeySize {
				t.Fatalf("expected key size error for overflow number, got %s", err)
			}
		}
	})
}
