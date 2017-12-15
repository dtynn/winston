package test

import (
	"reflect"
	"testing"

	"github.com/dtynn/winston/storage"
)

// Batch batch operations
func Batch(t *testing.T, s storage.Storage) {
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

	t.Run("BatchPut", func(t *testing.T) {
		batch, err := s.Batch()
		if err != nil {
			t.Errorf("get batch %s", err)
			return
		}

		defer batch.Close()

		for i, c := range cases {
			if err := batch.Put(c.key, c.val1); err != nil {
				t.Errorf("#%d put: %s", i+1, err)
				return
			}
		}

		if err := batch.Commit(); err != nil {
			t.Errorf("commit %s", err)
			return
		}

		for i, c := range cases {
			val, err := s.Get(c.key)
			if err != nil {
				t.Errorf("#%d got: %s", i+1, err)
				return
			}

			if !reflect.DeepEqual(val, c.val1) {
				t.Errorf("#%d expected %s, got %s", i+1, string(c.val1), string(val))
				return
			}
		}
	})

	t.Run("BatchUpdate", func(t *testing.T) {
		batch, err := s.Batch()
		if err != nil {
			t.Errorf("get batch %s", err)
			return
		}

		defer batch.Close()

		for i, c := range cases {
			if c.val2 == nil {
				if err := batch.Del(c.key); err != nil {
					t.Errorf("#%d del: %s", i+1, err)
					return
				}
				continue
			}

			if err := batch.Put(c.key, c.val2); err != nil {
				t.Errorf("#%d update: %s", i+1, err)
				return
			}
		}

		if err := batch.Commit(); err != nil {
			t.Errorf("commit %s", err)
			return
		}

		for i, c := range cases {
			val, err := s.Get(c.key)
			if err != nil {
				t.Errorf("#%d got: %s", i+1, err)
				return
			}

			if c.val2 == nil {
				if val != nil {
					t.Errorf("#%d expected nil, got %s", i+1, string(val))
				}

				continue
			}

			if !reflect.DeepEqual(val, c.val2) {
				t.Errorf("#%d expected updated value %s, got %s", i+1, string(c.val2), string(val))
				return
			}
		}
	})
}
