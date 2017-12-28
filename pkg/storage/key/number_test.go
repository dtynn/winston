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

	t.Run("Fixed", func(t *testing.T) {
		prefix := []byte("_fixed_")

		f1 := UI64(rand.Uint64())
		f2 := I64(rand.Int63())
		f3 := UI32(rand.Uint32())
		f4 := I32(rand.Int31())

		fs := FixedFormatters{&f1, &f2, &f3, &f4}
		key := Key(prefix, fs)

		var r1 UI64
		var r2 I64
		var r3 UI32
		var r4 I32

		rs := FixedFormatters{&r1, &r2, &r3, &r4}
		if err := Unmarshal(key, prefix, rs); err != nil {
			t.Fatal(err)
		}

		if r1 != f1 {
			t.Errorf("expected uint64 %d, got %d", f1, r1)
		}

		if r2 != f2 {
			t.Errorf("expected int64 %d, got %d", f2, r2)
		}

		if r3 != f3 {
			t.Errorf("expected uint32 %d, got %d", f3, r3)
		}

		if r4 != f4 {
			t.Errorf("expected int32 %d, got %d", f4, r4)
		}
	})
}
