package memstore

import (
	"fmt"
	"os"
	"testing"

	ds "github.com/ipfs/go-datastore"
)

var dataStore ds.Datastore
var key = ds.NewKey("foo")
var val = 101

func TestMemstoreBasicOps(t *testing.T) {
	err := dataStore.Put(key, val)
	if err != nil {
		t.Errorf("Could not store key")
	}

	fetchedVal, err := dataStore.Get(key)
	if err != nil || fetchedVal != val {
		t.Errorf("Error on Get, expected %d got %d", val, fetchedVal)
	}

	exists, err := dataStore.Has(key)
	if err != nil || !exists {
		t.Errorf("Error on Has, expected %v got %v", true, exists)
	}

	err = dataStore.Delete(key)
	if err != nil {
		t.Errorf("Could not delete key %v", key)
	}

	_, err = dataStore.Get(key)
	if err == nil {
		t.Errorf("Expected error on fetching deleted key")
	}

	exists, err = dataStore.Has(key)
	if err != nil || exists {
		t.Errorf("Expected Has to return %v (got %v)", false, exists)
	}
}

func TestMultiKeys(t *testing.T) {
	store := NewIntMemstore()
	expectedKeys := make([]ds.Key, 10)

	makeKey := func(k int) ds.Key {
		return ds.NewKey(fmt.Sprintf("key-%03d", k))
	}

	for k := 0; k < 10; k++ {
		key := makeKey(k)
		store.Put(key, k)
		expectedKeys[k] = key
	}

	for k, key := range store.Keys() {
		if expectedKey := makeKey(k); key != expectedKey {
			t.Errorf("Got unexpected key: %v (expected %v)", key, expectedKey)
		}
	}
}

func TestQuery(t *testing.T) {
	dataStore.Put(key, val)

	// query.Query()

	// dataStore.Query()
}

func TestMain(m *testing.M) {
	dataStore = NewIntMemstore()

	os.Exit(m.Run())
}
