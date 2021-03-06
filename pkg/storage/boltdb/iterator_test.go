package boltdb

import (
	"testing"

	"github.com/dtynn/winston/pkg/storage/test"
)

func TestBoltdbIterator(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.Iterator(t, s)
}
