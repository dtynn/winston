package test

import (
	"reflect"
	"testing"

	"github.com/dtynn/winston/pkg/storage"
)

// StorageUpdate test case for key update
func StorageUpdate(t *testing.T, s storage.Storage) {
	cases := []struct {
		key  []byte
		val1 []byte
		val2 []byte
	}{
		{
			key:  []byte("a"),
			val1: []byte("vala1"),
			val2: []byte("vala2"),
		},
		{
			key:  []byte("b"),
			val1: []byte("valb1"),
			val2: nil,
		},
		{
			key:  []byte("c"),
			val1: []byte("valc1"),
			val2: []byte("valc2"),
		},
		{
			key:  []byte("d"),
			val1: []byte("vald1"),
			val2: nil,
		},
		{
			key:  []byte("e"),
			val1: []byte("vale1"),
			val2: []byte("vale2"),
		},
	}

	mkeys := make([][]byte, len(cases))

	t.Run("Put", func(t *testing.T) {
		for i, c := range cases {
			if err := s.Put(c.key, c.val1); err != nil {
				t.Errorf("#%d put val1: %s", i+1, err)
				return
			}

			mkeys[i] = c.key
		}

		for i, c := range cases {
			got, err := s.Get(c.key)
			if err != nil {
				t.Errorf("#%d get val1: %s", i+1, err)
				return
			}

			if !reflect.DeepEqual(got, c.val1) {
				t.Errorf("#%d expected val1 %s, got %s", i+1, string(c.val1), string(got))
			}
		}
	})

	t.Run("UpdateAndDel", func(t *testing.T) {
		for i, c := range cases {
			if c.val2 != nil {
				if err := s.Put(c.key, c.val2); err != nil {
					t.Errorf("#%d put val2: %s", i+1, err)
					return
				}
				continue
			}

			if err := s.Del(c.key); err != nil {
				t.Errorf("#%d del: %s", i+1, err)
				return
			}
		}
	})

	t.Run("MGet", func(t *testing.T) {
		vals, err := s.MGet(mkeys...)
		if err != nil {
			t.Errorf("mget %s", err)
			return
		}

		if len(vals) != len(cases) {
			t.Errorf("expected %d values, got %d", len(cases), len(vals))
			return
		}

		for i, v := range vals {
			if !reflect.DeepEqual(v, cases[i].val2) {
				t.Errorf("#%d expected %v, got %v", i+1, cases[i].val2, v)
				return
			}
		}
	})
}
