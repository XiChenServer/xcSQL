package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/logs"
	"fmt"
	"time"
)

const (
	String    uint16 = 1
	StringSet uint16 = 2
)

func (db XcDB) Set(key, value []byte, ttl ...time.Duration) error {
	err := db.doSet(key, value, ttl...)
	if err != nil {
		return err
	}
	return nil
}

func (db *XcDB) doSet(key, value []byte, ttl ...time.Duration) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	e := NewKeyValueEntry(key, value, String, StringSet, ttl...)
	stroeLocal, err := db.storageManager.StoreData(e)
	if err != nil {
		logs.SugarLogger.Error("string set fail:", err)
		return err
	}
	datainfo := &lsm.DataInfo{
		DataMeta:        *e.DataMeta,
		StorageLocation: stroeLocal,
	}
	db.lsm.Insert(key, datainfo)
	fmt.Println(stroeLocal.Size, stroeLocal.Offset, string(stroeLocal.FileName))
	return nil
}
func (db *XcDB) Get(key []byte) (*model.KeyValue, error) {
	data, err := db.doGet(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (db *XcDB) doGet(key []byte) (*model.KeyValue, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	datainfo, err := db.lsm.Get(key)
	if err != nil {
	}
	offset := datainfo.Offset
	fileName := datainfo.FileName
	size := datainfo.Size
	data, err := db.storageManager.DecompressAndFillData(string(fileName), offset, size)
	if err != nil {
		return nil, err
	}
	return data, nil
}
