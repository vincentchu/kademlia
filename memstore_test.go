package memstore

import (
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

func TestQuery(t *testing.T) {
	dataStore.Put(key, val)

	// query.Query()

	// dataStore.Query()
}

func TestMain(m *testing.M) {
	dataStore = NewIntMemstore(100)

	os.Exit(m.Run())
}
