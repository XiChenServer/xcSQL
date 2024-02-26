package database

import (
	"SQL/internal/storage"
	"time"
)

const (
	String    uint16 = 1
	StringSet uint16 = 2
)

func (db XcDB) Set(key, value []byte, ttl ...time.Duration) {

}

func (db *XcDB) doSet(key, value []byte, ttl ...time.Duration) {
	db.mu.Lock()
	defer db.mu.Unlock()
	e := NewKeyValueEntry(key, value, String, StringSet)
	local, err := storage.StorageManager.StoreData(e)
	if err != nil {

	}
}
