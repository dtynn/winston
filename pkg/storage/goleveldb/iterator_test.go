package goleveldb

import (
	"testing"

	"github.com/dtynn/winston/pkg/storage/test"
)

func TestGoLeveldbIterator(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.Iterator(t, s)
}
