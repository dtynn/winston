package boltdb

import (
	"testing"

	"github.com/coreos/bbolt"
	"github.com/dtynn/winston/storage/test"
)

func setupTestStorage(t *testing.T) *Storage {
	s, err := Open("./testdb")
	if err != nil {
		t.Fatalf("open storage: %s", err)
	}

	return s
}

func teardownTestStorage(s *Storage) {
	s.Close()
	s.cleanup()
}

func TestBoltdbOption(t *testing.T) {
	t.Run("OptionBucket", func(t *testing.T) {
		bname := "_testbucket"
		s, err := Open("./testdb", Bucket([]byte(bname)))
		if err != nil {
			t.Fatalf("open db: %s", err)
		}

		defer teardownTestStorage(s)

		s.db.View(func(tx *bolt.Tx) error {
			if b := tx.Bucket(defaultBucket); b != nil {
				t.Errorf("expcted nil default bucket")
			}

			if b := tx.Bucket([]byte(bname)); b == nil {
				t.Errorf("expected specified bucket, got nil")
			}

			return nil
		})
	})
}

func TestBoltdbPut(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.StoragePut(t, s)
}

func TestBoltdbUpdate(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.StorageUpdate(t, s)
}

func TestBoltdbDel(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.StorageDel(t, s)
}
