package storage

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	_testStorage *KeyValueInMem
	_testKey     = "test-key"
	_testValue   = []byte("test-value")
	_testExp     = 3 * time.Second
)

func TestMain(m *testing.M) {
	_testStorage = NewKeyValueInMem()
	os.Exit(m.Run())
}

func TestKeyValueInMem_Set(t *testing.T) {
	// set
	t.Logf("Set %s (key) %s (value) expired at %s", _testKey, _testValue, _testExp.String())
	err := _testStorage.Set(_testKey, _testValue, _testExp)
	require.NoError(t, err)
}

func TestKeyValueInMem_Get(t *testing.T) {
	// get
	t.Logf("Get %s value", _testKey)
	val, err := _testStorage.Get(_testKey)
	require.NoError(t, err)

	// log result
	t.Logf("%s value: %s", _testKey, val)
	require.Equal(t, _testValue, val)

	// wait for expiration
	wait := _testExp + time.Second
	t.Logf("Wait %s for key-value expiration...", wait.String())
	time.Sleep(wait)

	// get after expiration
	t.Logf("Get %s value", _testKey)
	_, err = _testStorage.Get(_testKey)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestKeyValueInMem_Delete(t *testing.T) {
	// set
	t.Logf("Set %s (key) %s (value) expired at %s", _testKey, _testValue, _testExp.String())
	err := _testStorage.Set(_testKey, _testValue, _testExp)
	require.NoError(t, err)

	// delete
	t.Logf("Delete %s (key) %s (value)", _testKey, _testValue)
	err = _testStorage.Delete(_testKey)
	require.NoError(t, err)

	// get after deleting
	t.Logf("Get %s value", _testKey)
	_, err = _testStorage.Get(_testKey)
	require.ErrorIs(t, err, ErrNotFound)
}
