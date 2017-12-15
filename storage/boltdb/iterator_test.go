package boltdb

import (
	"testing"

	"github.com/dtynn/winston/storage/test"
)

func TestBoltdbIteratorNormal(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.IteratorNormal(t, s)
}

func TestBoltdbIteratorPrefix(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.IteratorPrefix(t, s)
}

func TestBoltdbIteratorRange(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.IteratorRange(t, s)
}
