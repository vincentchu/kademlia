package memstore

import (
	"fmt"
	"sync"

	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
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

	// _, exists := store.memMap[key]
	// if exists {

	// } else {
	// 	store.keys
	// }

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

// Delete deletes stuff
func (store *IntMemstore) Delete(key ds.Key) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	delete(store.memMap, key)

	return nil
}

// Query queries stuff
func (store *IntMemstore) Query(q query.Query) (query.Results, error) {
	// rb := query.NewResultBuilder(q)
	// return rb.Results(), nil

	// return query.Results{}, nil
	panic(1)
}

// NewIntMemstore creates new Integer Memstore
func NewIntMemstore(capacity int) *IntMemstore {
	return &IntMemstore{
		sync.Mutex{},
		make(map[ds.Key]int, capacity),
		make([]ds.Key, capacity),
	}
}
