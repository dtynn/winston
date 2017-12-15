package goleveldb

import (
	"testing"

	"github.com/dtynn/winston/storage/test"
)

func TestGoLeveldbBatch(t *testing.T) {
	s := setupTestStorage(t)
	defer teardownTestStorage(s)

	test.Batch(t, s)
}
