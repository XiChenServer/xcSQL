package database

import (
	"SQL/internal/lsm"
	"SQL/logs"
	"fmt"
	"time"
)

const (
	String    uint16 = 1
	StringSet uint16 = 2
)

func (db XcDB) Set(key, value []byte, ttl ...time.Duration) {
	db.doSet(key, value, ttl...)
}

func (db *XcDB) doSet(key, value []byte, ttl ...time.Duration) {
	db.mu.Lock()
	defer db.mu.Unlock()
	e := NewKeyValueEntry(key, value, String, StringSet, ttl...)
	stroeLocal, err := db.storageManager.StoreData(e)
	if err != nil {
		logs.SugarLogger.Error("string set fail:", err)
	}
	datainfo := &lsm.DataInfo{
		DataMeta:        *e.DataMeta,
		StorageLocation: stroeLocal,
	}
	db.lsm.Insert(key, datainfo)
	fmt.Println(stroeLocal.Size, stroeLocal.Offset, string(stroeLocal.FileName))
}
