package goleveldb

import (
	"testing"

	"github.com/dtynn/winston/storage/test"
)

func TestGoLeveldbIteratorNormal(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.IteratorNormal(t, s)
}

func TestGoLeveldbIteratorPrefix(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.IteratorPrefix(t, s)
}

func TestGoLeveldbIteratorRange(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.IteratorRange(t, s)
}
