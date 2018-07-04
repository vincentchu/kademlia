package memstore

import (
	"fmt"
	"os"
	"testing"

	ds "gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore"
	query "gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore/query"
)

var dataStore ds.Datastore
var key = ds.NewKey("foo")
var val = 101

func makeKey(k int) ds.Key {
	return ds.NewKey(fmt.Sprintf("key-%03d", k))
}

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
	for k := 0; k < 10; k++ {
		dataStore.Put(makeKey(k), k)
	}

	q := query.Query{
		Prefix:   "",
		Filters:  nil,
		Orders:   nil,
		Limit:    -1,
		Offset:   -1,
		KeysOnly: false,
	}

	results, err := dataStore.Query(q)
	if err != nil {
		t.Errorf("Unexpected err %v encountered with querying", err)
	}

	ctr := 0
	for result := range results.Next() {
		if keyStr := makeKey(ctr).String(); result.Key != keyStr {
			t.Errorf("Unexpected key: %s (expected %s)", result.Key, keyStr)
		}

		if result.Value != ctr {
			t.Errorf("Unexpected value: %d (expected %d)", result.Value, ctr)
		}

		ctr++
	}
}

func TestMain(m *testing.M) {
	dataStore = NewIntMemstore()

	os.Exit(m.Run())
}
