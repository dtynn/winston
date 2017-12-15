package test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/dtynn/winston/storage"
)

// IteratorNormal normal iterator
func IteratorNormal(t *testing.T, s storage.Storage) {
	keys := []string{
		"a",
		"b",
		"c",
		"d",
		"e",
	}

	val := make([]byte, 256)
	if _, err := rand.Read(val); err != nil {
		t.Errorf("rand.Read %s", err)
		return
	}

	for i, k := range keys {
		if err := s.Put([]byte(k), val); err != nil {
			t.Errorf("#%d put: %s", i+1, err)
			return
		}
	}

	iter, err := s.PrefixIterator(nil)
	if err != nil {
		t.Errorf("normal iterator %s", err)
		return
	}

	defer iter.Close()

	got := make([]string, 0, len(keys))
	for iter.First(); iter.Valid(); iter.Next() {
		k := iter.Key()
		got = append(got, string(k))
	}

	if !reflect.DeepEqual(keys, got) {
		t.Errorf("expected keys %v, got %v", keys, got)
		return
	}
}

// IteratorPrefix normal iterator
func IteratorPrefix(t *testing.T, s storage.Storage) {
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

	testFn := func(t *testing.T, prefix []byte, expected []string) {
		iter, err := s.PrefixIterator(prefix)
		if err != nil {
			t.Errorf("prefix iterator for %q %s", string(prefix), err)
			return
		}

		defer iter.Close()

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

	t.Run("IteratorPrefix.Prefix_nil", func(t *testing.T) {
		testFn(t, nil, keys)
	})

	t.Run("IteratorPrefix.Prefix_a", func(t *testing.T) {
		testFn(t, []byte("a"), []string{"a", "a1", "a2", "a3"})
	})

	t.Run("IteratorPrefix.Prefix_c", func(t *testing.T) {
		testFn(t, []byte("c"), []string{"c", "c1", "c2"})
	})

	t.Run("IteratorPrefix.Prefix_d", func(t *testing.T) {
		testFn(t, []byte("d"), []string{"d1", "d2"})
	})

	t.Run("IteratorPrefix.Prefix_e", func(t *testing.T) {
		testFn(t, []byte("e"), []string{})
	})
}

// IteratorRange normal iterator
func IteratorRange(t *testing.T, s storage.Storage) {
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

	t.Run("IteratorRange.Put", func(t *testing.T) {
		for i, k := range keys {
			if err := s.Put([]byte(k), val); err != nil {
				t.Errorf("#%d put: %s", i+1, err)
				return
			}
		}
	})

	testFn := func(t *testing.T, start, end []byte, expected []string) {
		iter, err := s.RangeIterator(start, end)
		if err != nil {
			t.Errorf("range iterator for [%q, %q) %s", string(start), string(end), err)
			return
		}

		defer iter.Close()

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

	t.Run("IteratorRange.Range_nil_nil", func(t *testing.T) {
		testFn(t, nil, nil, keys)
	})

	t.Run("IteratorRange.Range_nil_c", func(t *testing.T) {
		testFn(t, nil, []byte("c"), []string{"a", "a1", "a2", "a3", "b"})
	})

	t.Run("IteratorRange.Range_c_nil", func(t *testing.T) {
		testFn(t, []byte("c"), nil, []string{"c", "c1", "c2", "d1", "d2", "f", "g"})
	})

	t.Run("IteratorRange.Range_d_f", func(t *testing.T) {
		testFn(t, []byte("d"), []byte("f"), []string{"d1", "d2"})
	})

	t.Run("IteratorRange.Range_x_z", func(t *testing.T) {
		testFn(t, []byte("x"), []byte("z"), []string{})
	})
}
