package test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/dtynn/winston/pkg/storage"
)

// Iterator normal iterator
func Iterator(t *testing.T, s storage.Storage) {
	keys := []string{
		"a",
		"a1",
		"a2",
		"a3",
		"b",
		"c",
		"c1",
		"c2",
		"d1",
		"d2",
		"f",
		"g",
	}

	val := make([]byte, 256)
	rand.Read(val)

	t.Run("IteratorPrefix.Put", func(t *testing.T) {
		for i, k := range keys {
			if err := s.Put([]byte(k), val); err != nil {
				t.Errorf("#%d put: %s", i+1, err)
				return
			}
		}
	})

	testPrefixFn := func(t *testing.T, prefix []byte, expected []string) {
		iter, err := s.PrefixIterator(prefix)
		if err != nil {
			t.Errorf("prefix iterator for %q %s", string(prefix), err)
			return
		}

		defer iter.Close()

		{
			got := make([]string, 0, len(expected))
			for iter.First(); iter.Valid(); iter.Next() {
				k, v := iter.Key(), iter.Value()
				if !reflect.DeepEqual(v, val) {
					t.Errorf("unexpected value during iteration")
					return
				}

				got = append(got, string(k))
			}

			if err := iter.Err(); err != nil {
				t.Errorf("got iter err: %s", err)
				return
			}

			if !reflect.DeepEqual(expected, got) {
				t.Errorf("expected keys %v, got %v", expected, got)
				return
			}
		}

		// reverse
		{
			got := make([]string, 0, len(expected))
			for iter.Last(); iter.Valid(); iter.Prev() {
				k, v := iter.Key(), iter.Value()
				if !reflect.DeepEqual(v, val) {
					t.Errorf("unexpected value during iteration")
					return
				}

				got = append([]string{string(k)}, got...)
			}

			if err := iter.Err(); err != nil {
				t.Errorf("got iter err: %s", err)
				return
			}

			if !reflect.DeepEqual(expected, got) {
				t.Errorf("expected keys %v, got %v", expected, got)
				return
			}
		}
	}

	testPrefixFnWhileNext := func(t *testing.T, prefix []byte, expected []string) {
		iter, err := s.PrefixIterator(prefix)
		if err != nil {
			t.Errorf("prefix iterator for %q %s", string(prefix), err)
			return
		}

		defer iter.Close()

		got := make([]string, 0, len(expected))
		for iter.Next() {
			k, v := iter.Key(), iter.Value()
			if !reflect.DeepEqual(v, val) {
				t.Errorf("unexpected value during iteration")
				return
			}

			got = append(got, string(k))
		}

		if err := iter.Err(); err != nil {
			t.Errorf("got iter err: %s", err)
			return
		}

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected keys %v, got %v", expected, got)
			return
		}
	}

	testRangeFn := func(t *testing.T, start, end []byte, expected []string) {
		iter, err := s.RangeIterator(start, end)
		if err != nil {
			t.Errorf("range iterator for [%q, %q) %s", string(start), string(end), err)
			return
		}

		defer iter.Close()

		{
			got := make([]string, 0, len(expected))
			for iter.First(); iter.Valid(); iter.Next() {
				k, v := iter.Key(), iter.Value()
				if !reflect.DeepEqual(v, val) {
					t.Errorf("unexpected value during iteration")
					return
				}

				got = append(got, string(k))
			}

			if err := iter.Err(); err != nil {
				t.Errorf("got iter err: %s", err)
				return
			}

			if !reflect.DeepEqual(expected, got) {
				t.Errorf("expected keys %v, got %v", expected, got)
				return
			}
		}

		// reverse
		{
			got := make([]string, 0, len(expected))
			for iter.Last(); iter.Valid(); iter.Prev() {
				k, v := iter.Key(), iter.Value()
				if !reflect.DeepEqual(v, val) {
					t.Errorf("unexpected value during iteration")
					return
				}

				got = append([]string{string(k)}, got...)
			}

			if err := iter.Err(); err != nil {
				t.Errorf("got iter err: %s", err)
				return
			}

			if !reflect.DeepEqual(expected, got) {
				t.Errorf("expected keys %v, got %v", expected, got)
				return
			}
		}
	}

	testRangeFnWhileNext := func(t *testing.T, start, end []byte, expected []string) {
		iter, err := s.RangeIterator(start, end)
		if err != nil {
			t.Errorf("range iterator for [%q, %q) %s", string(start), string(end), err)
			return
		}

		defer iter.Close()

		got := make([]string, 0, len(expected))
		for iter.Next() {
			k, v := iter.Key(), iter.Value()
			if !reflect.DeepEqual(v, val) {
				t.Errorf("unexpected value during iteration")
				return
			}

			got = append(got, string(k))
		}

		if err := iter.Err(); err != nil {
			t.Errorf("got iter err: %s", err)
			return
		}

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected keys %v, got %v", expected, got)
			return
		}
	}

	t.Run("IteratorPrefix.Prefix_nil", func(t *testing.T) {
		testPrefixFn(t, nil, keys)
		testPrefixFnWhileNext(t, nil, keys)
	})

	t.Run("IteratorPrefix.Prefix_a", func(t *testing.T) {
		testPrefixFn(t, []byte("a"), []string{"a", "a1", "a2", "a3"})
		testPrefixFnWhileNext(t, []byte("a"), []string{"a", "a1", "a2", "a3"})
	})

	t.Run("IteratorPrefix.Prefix_c", func(t *testing.T) {
		testPrefixFn(t, []byte("c"), []string{"c", "c1", "c2"})
		testPrefixFnWhileNext(t, []byte("c"), []string{"c", "c1", "c2"})
	})

	t.Run("IteratorPrefix.Prefix_d", func(t *testing.T) {
		testPrefixFn(t, []byte("d"), []string{"d1", "d2"})
		testPrefixFnWhileNext(t, []byte("d"), []string{"d1", "d2"})
	})

	t.Run("IteratorPrefix.Prefix_e", func(t *testing.T) {
		testPrefixFn(t, []byte("e"), []string{})
		testPrefixFnWhileNext(t, []byte("e"), []string{})
	})

	t.Run("IteratorRange.Range_nil_nil", func(t *testing.T) {
		testRangeFn(t, nil, nil, keys)
		testRangeFnWhileNext(t, nil, nil, keys)
	})

	t.Run("IteratorRange.Range_nil_c", func(t *testing.T) {
		testRangeFn(t, nil, []byte("c"), []string{"a", "a1", "a2", "a3", "b"})
		testRangeFnWhileNext(t, nil, []byte("c"), []string{"a", "a1", "a2", "a3", "b"})
	})

	t.Run("IteratorRange.Range_c_nil", func(t *testing.T) {
		testRangeFn(t, []byte("c"), nil, []string{"c", "c1", "c2", "d1", "d2", "f", "g"})
		testRangeFnWhileNext(t, []byte("c"), nil, []string{"c", "c1", "c2", "d1", "d2", "f", "g"})
	})

	t.Run("IteratorRange.Range_d_f", func(t *testing.T) {
		testRangeFn(t, []byte("d"), []byte("f"), []string{"d1", "d2"})
		testRangeFnWhileNext(t, []byte("d"), []byte("f"), []string{"d1", "d2"})
	})

	t.Run("IteratorRange.Range_x_z", func(t *testing.T) {
		testRangeFn(t, []byte("x"), []byte("z"), []string{})
		testRangeFnWhileNext(t, []byte("x"), []byte("z"), []string{})
	})
}
