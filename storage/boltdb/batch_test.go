package boltdb

import (
	"testing"

	"github.com/dtynn/winston/storage/test"
)

func TestBoltdbBatch(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.Batch(t, s)
}
