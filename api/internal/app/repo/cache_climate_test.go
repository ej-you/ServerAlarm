package repo

import (
	"os"
	"testing"
	"time"

	"server-alarm/api/internal/app/entity"
	"server-alarm/api/internal/pkg/storage"

	"github.com/stretchr/testify/require"
)

var _testClimateRepoCache *ClimateRepoCache

func TestMain(m *testing.M) {
	storageInst := storage.NewKeyValueInMem()
	historyLen := 2

	_testClimateRepoCache = NewClimateRepoCache(storageInst, historyLen)

	os.Exit(m.Run())
}

func TestHistory(t *testing.T) {
	t.Log("Get empty history")
	history, err := _testClimateRepoCache.GetHistory()
	require.NoError(t, err, "get empty history")
	t.Logf("History: %+v", history)

	t.Log("Set 1st climate to history and get new history")
	err = _testClimateRepoCache.SetHistory(&entity.Climate{
		ID:          1,
		Datetime:    time.Now(),
		Temperature: 29.8,
	})
	require.NoError(t, err, "set 1st climate to history")
	printHistory(t)

	t.Log("Set 2nd climate to history and get new history")
	err = _testClimateRepoCache.SetHistory(&entity.Climate{
		ID:          2,
		Datetime:    time.Now(),
		Temperature: 30.0,
	})
	require.NoError(t, err, "set 2nd climate to history")
	printHistory(t)

	t.Log("Set 3st climate to history and get new history")
	err = _testClimateRepoCache.SetHistory(&entity.Climate{
		ID:          3,
		Datetime:    time.Now(),
		Temperature: 29.9,
	})
	require.NoError(t, err, "set 3st climate to history")
	printHistory(t)
}

func printHistory(t *testing.T) {
	history, err := _testClimateRepoCache.GetHistory()
	require.NoError(t, err, "get history")
	t.Logf("History (len: %d | cap: %d):", len(history), cap(history))
	for _, elem := range history {
		t.Logf("\t%+v", elem)
	}
}
