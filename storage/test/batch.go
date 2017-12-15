package test

import (
	"testing"

	"github.com/dtynn/winston/storage"
)

// Batch batch operations
func Batch(t *testing.T, s storage.Storage) {
	cases := []struct {
		key  string
		val1 string
		val2 string
	}{
		{
			key:  "a",
			val1: "vala1",
			val2: "vala2",
		},
		{
			key:  "b",
			val1: "valb1",
			val2: "",
		},
		{
			key:  "c",
			val1: "valc1",
			val2: "valc2",
		},
		{
			key:  "d",
			val1: "vald1",
			val2: "",
		},
		{
			key:  "e",
			val1: "vale1",
			val2: "vale2",
		},
	}

	t.Run("BatchInit", func(t *testing.T) {
		batch, err := s.Batch()
		if err != nil {
			t.Errorf("get batch %s", err)
			return
		}

		defer batch.Close()

		for i, c := range cases {
			if err := batch.Put([]byte(c.key), []byte(c.val1)); err != nil {
				t.Errorf("#%d put: %s", i+1, err)
				return
			}
		}

		if err := batch.Commit(); err != nil {
			t.Errorf("commit %s", err)
			return
		}

		for i, c := range cases {
			val, err := s.Get([]byte(c.key))
			if err != nil {
				t.Errorf("#%d got: %s", i+1, err)
				return
			}

			if s := string(val); s != c.val1 {
				t.Errorf("#%d expected %s, got %s", i+1, c.val1, s)
				return
			}
		}
	})

	t.Run("BatchInit", func(t *testing.T) {
		batch, err := s.Batch()
		if err != nil {
			t.Errorf("get batch %s", err)
			return
		}

		defer batch.Close()

		for i, c := range cases {
			if c.val2 == "" {
				if err := batch.Del([]byte(c.key)); err != nil {
					t.Errorf("#%d del: %s", i+1, err)
					return
				}
				continue
			}

			if err := batch.Put([]byte(c.key), []byte(c.val2)); err != nil {
				t.Errorf("#%d update: %s", i+1, err)
				return
			}
		}

		if err := batch.Commit(); err != nil {
			t.Errorf("commit %s", err)
			return
		}

		for i, c := range cases {
			val, err := s.Get([]byte(c.key))
			if err != nil {
				t.Errorf("#%d got: %s", i+1, err)
				return
			}

			if c.val2 == "" {
				if val != nil {
					t.Errorf("#%d expected nil, got %s", i+1, string(val))
				}

				continue
			}

			if s := string(val); s != c.val2 {
				t.Errorf("#%d expected updated value %s, got %s", i+1, c.val2, s)
				return
			}
		}
	})
}
