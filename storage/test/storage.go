package test

import (
	"testing"

	"github.com/dtynn/winston/storage"
)

// StoragePut test case for key put
func StoragePut(t *testing.T, s storage.Storage) {
	cases := []struct {
		key string
		val string
	}{
		{
			key: "a",
			val: "vala",
		},
		{
			key: "b",
			val: "valb",
		},
		{
			key: "c",
			val: "valc",
		},
		{
			key: "d",
			val: "vald",
		},
		{
			key: "e",
			val: "vale",
		},
	}

	for i, c := range cases {
		if err := s.Put([]byte(c.key), []byte(c.val)); err != nil {
			t.Errorf("#%d put: %s", i+1, err)
			return
		}
	}

	for i, c := range cases {
		got, err := s.Get([]byte(c.key))
		if err != nil {
			t.Errorf("#%d get: %s", i+1, err)
			return
		}

		if s := string(got); s != c.val {
			t.Errorf("#%d expected %s, got %s", i+1, c.val, s)
		}
	}
}

// StorageUpdate test case for key update
func StorageUpdate(t *testing.T, s storage.Storage) {
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
			val2: "valb2",
		},
		{
			key:  "c",
			val1: "valc1",
			val2: "valc2",
		},
		{
			key:  "d",
			val1: "vald1",
			val2: "vald2",
		},
		{
			key:  "e",
			val1: "vale1",
			val2: "vale2",
		},
	}

	for i, c := range cases {
		if err := s.Put([]byte(c.key), []byte(c.val1)); err != nil {
			t.Errorf("#%d put val1: %s", i+1, err)
			return
		}
	}

	for i, c := range cases {
		got, err := s.Get([]byte(c.key))
		if err != nil {
			t.Errorf("#%d get val1: %s", i+1, err)
			return
		}

		if s := string(got); s != c.val1 {
			t.Errorf("#%d expected val1 %s, got %s", i+1, c.val1, s)
		}
	}

	for i, c := range cases {
		if err := s.Put([]byte(c.key), []byte(c.val2)); err != nil {
			t.Errorf("#%d put val2: %s", i+1, err)
			return
		}
	}

	for i, c := range cases {
		got, err := s.Get([]byte(c.key))
		if err != nil {
			t.Errorf("#%d get val2: %s", i+1, err)
			return
		}

		if s := string(got); s != c.val2 {
			t.Errorf("#%d expected val2 %s, got %s", i+1, c.val2, s)
		}
	}
}

// StorageDel test case for key del
func StorageDel(t *testing.T, s storage.Storage) {
	cases := []struct {
		key string
		val string
	}{
		{
			key: "a",
			val: "vala",
		},
		{
			key: "b",
			val: "valb",
		},
		{
			key: "c",
			val: "valc",
		},
		{
			key: "d",
			val: "vald",
		},
		{
			key: "e",
			val: "vale",
		},
	}

	for i, c := range cases {
		if err := s.Put([]byte(c.key), []byte(c.val)); err != nil {
			t.Errorf("#%d put: %s", i+1, err)
			return
		}
	}

	for i, c := range cases {
		got, err := s.Get([]byte(c.key))
		if err != nil {
			t.Errorf("#%d get: %s", i+1, err)
			return
		}

		if s := string(got); s != c.val {
			t.Errorf("#%d expected %s, got %s", i+1, c.val, s)
		}
	}

	for i, c := range cases {
		if err := s.Del([]byte(c.key)); err != nil {
			t.Errorf("#%d del: %s", i+1, err)
			return
		}
	}

	for i, c := range cases {
		got, err := s.Get([]byte(c.key))
		if err != nil {
			t.Errorf("#%d get after del: %s", i+1, err)
			return
		}

		if got != nil {
			t.Errorf("#%d expected nil, got %s", i+1, s)
		}
	}
}
