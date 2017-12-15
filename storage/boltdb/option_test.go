package boltdb

import (
	"testing"

	"github.com/coreos/bbolt"
)

func TestBoltdbOption(t *testing.T) {
	t.Run("OptionBucket", func(t *testing.T) {
		bname := "_testbucket"
		s, err := Open("./testdb/test.db", Bucket([]byte(bname)))
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
