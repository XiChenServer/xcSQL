package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/logs"
	"errors"
	"fmt"
	"time"
)

const (
	String    uint16 = 1
	StringSet uint16 = 2
)

func (db XcDB) Set(key, value []byte, ttl ...uint64) error {
	err := db.doSet(key, value, ttl...)
	if err != nil {
		return err
	}
	return nil
}

func (db *XcDB) doSet(key, value []byte, ttl ...uint64) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	var timeSlice []time.Duration
	for _, t := range ttl {
		timeSlice = append(timeSlice, time.Duration(t)*time.Second)
	}

	e := NewKeyValueEntry(key, value, String, StringSet, timeSlice...)
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
		fmt.Println(err)
		return nil, err
	}

	if isExpired(data) {
		err = errors.New("Data has expired")
		return nil, err
	}

	err = db.reSet(data)
	if err != nil {
		return nil, err
	}

	if data == nil {
		err := errors.New("No data found")
		return nil, err
	}
	return data, nil
}

func (db *XcDB) reSet(data *model.KeyValue) error {
	now := time.Now()
	data.AccessTime = now
	stroeLocal, err := db.storageManager.StoreData(data)
	if err != nil {
		logs.SugarLogger.Error("string set fail:", err)
		return err
	}
	datainfo := &lsm.DataInfo{
		DataMeta:        *data.DataMeta,
		StorageLocation: stroeLocal,
	}
	err = db.lsm.Insert(data.DataMeta.Key, datainfo)
	if err != nil {
		return err
	}
	return nil
}

func isExpired(kv *model.KeyValue) bool {
	// 获取当前时间
	currentTime := time.Now()

	// 如果TTL为0，表示永不过期
	if kv.DataMeta.TTL == 0 {
		return false
	}
	// 计算过期时间
	expirationTime := kv.CreateTime.Add(kv.DataMeta.TTL)

	// 判断当前时间是否已经超过过期时间
	return currentTime.After(expirationTime)
}
