package memstore

import (
	"fmt"
	"sync"

	ds "gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore"
	query "gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore/query"
)

// IntMemstore does stuff
type IntMemstore struct {
	mutex  sync.Mutex
	memMap map[ds.Key]int
	keys   []ds.Key
}

// Put does stuff
func (store *IntMemstore) Put(key ds.Key, value interface{}) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	castValue, ok := value.(int)
	if !ok {
		return fmt.Errorf("Non-integer value detected: %v", value)
	}

	_, exists := store.memMap[key]
	if !exists {
		store.keys = append(store.keys, key)
	}

	store.memMap[key] = castValue

	return nil
}

// Get gets stuff
func (store *IntMemstore) Get(key ds.Key) (interface{}, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	value, exists := store.memMap[key]

	var err error
	if !exists {
		err = fmt.Errorf("Couldn't find value")
	}

	return value, err
}

// Has has stuff
func (store *IntMemstore) Has(key ds.Key) (bool, error) {
	_, err := store.Get(key)

	return err == nil, nil
}

// Keys list all currently stored Keys
func (store *IntMemstore) Keys() []ds.Key {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	return store.keys
}

// Commit is a stubbed out commit
func (store *IntMemstore) Commit() error {
	return nil
}

// Batch returns underlying store to fake batching behavior
func (store *IntMemstore) Batch() (ds.Batch, error) {
	return store, nil
}

// Delete deletes stuff
func (store *IntMemstore) Delete(key ds.Key) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, exists := store.memMap[key]

	if exists {
		for k := 0; k < len(store.keys); k++ {
			if store.keys[k] == key {
				store.keys = append(store.keys[:k], store.keys[k+1:]...)
				break
			}
		}

		delete(store.memMap, key)
	}

	return nil
}

// Query queries stuff
//
// For now, just return all results and can focus on the actual Query
// logic a bit later
func (store *IntMemstore) Query(q query.Query) (query.Results, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	entries := make([]query.Entry, len(store.keys))
	for idx, key := range store.keys {
		entries[idx] = query.Entry{
			Key:   key.String(),
			Value: store.memMap[key],
		}
	}

	return query.ResultsWithEntries(q, entries), nil
}

// NewIntMemstore creates new Integer Memstore
func NewIntMemstore() *IntMemstore {
	return &IntMemstore{
		sync.Mutex{},
		make(map[ds.Key]int),
		make([]ds.Key, 0),
	}
}
