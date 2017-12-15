package goleveldb

import (
	"testing"

	"github.com/dtynn/winston/storage/test"
)

func setupTestStorage(t *testing.T) *Storage {
	s, err := Open("./testdb/test.db")
	if err != nil {
		t.Fatalf("open storage: %s", err)
	}

	return s
}

func teardownTestStorage(s *Storage) {
	s.Close()
	s.cleanup()
}

func TestGoLeveldbUpdate(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.StorageUpdate(t, s)
}
