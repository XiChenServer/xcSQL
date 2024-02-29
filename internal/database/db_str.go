package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/logs"
	"errors"
	"fmt"
	"time"
)

const ()

// 对于字符串进行建立的操作
func (db XcDB) Set(key, value []byte, ttl ...uint64) error {
	err := db.doSet(key, value, ttl...)
	if err != nil {
		return err
	}
	return nil
}

func (db *XcDB) doSet(key, value []byte, ttl ...uint64) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	var timeSlice []time.Duration
	for _, t := range ttl {
		timeSlice = append(timeSlice, time.Duration(t)*time.Second)
	}

	e := NewKeyValueEntry(key, value, model.String, model.StringSet, timeSlice...)
	stroeLocal, err := db.StorageManager.StoreData(e)
	if err != nil {
		logs.SugarLogger.Error("string set fail:", err)
		return err
	}
	datainfo := &lsm.DataInfo{
		DataMeta:        *e.DataMeta,
		StorageLocation: stroeLocal,
	}
	lsmMap := *db.Lsm
	tree := lsmMap[model.String]
	tree.Insert(key, datainfo)
	return nil
}

// 获取字符串的操作
func (db *XcDB) Get(key []byte) (*model.KeyValue, error) {
	data, err := db.doGet(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (db *XcDB) doGet(key []byte) (*model.KeyValue, error) {
	db.Mu.RLock()
	defer db.Mu.RUnlock()
	lsmMap := *db.Lsm
	tree := lsmMap[model.String]
	datainfo, err := tree.Get(key)
	if err != nil {
		return nil, err
	}
	offset := datainfo.Offset
	fileName := datainfo.FileName
	size := datainfo.Size

	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)

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

// 访问之后，对于数据进行一定的修改，重新保存
func (db *XcDB) reSet(data *model.KeyValue) error {
	data.DataMeta.ValueSize = uint32(len(data.DataMeta.Value))
	now := time.Now()
	data.AccessTime = now

	stroeLocal, err := db.StorageManager.StoreData(data)
	if err != nil {
		logs.SugarLogger.Error("string set fail:", err)
		return err
	}
	datainfo := &lsm.DataInfo{
		DataMeta:        *data.DataMeta,
		StorageLocation: stroeLocal,
	}
	lsmMap := *db.Lsm
	tree := lsmMap[model.String]
	err = tree.Insert(data.DataMeta.Key, datainfo)
	if err != nil {
		return err
	}
	return nil
}

// 验证过期时间
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

// 获取数据的长度
func (db *XcDB) Strlen(key []byte) (uint32, error) {
	valueLen, err := db.doGetStrLen(key)
	if err != nil {
		return 0, err
	}
	if int(valueLen) == 0 {
		err = errors.New("Not found")
		return 0, err
	}
	return valueLen, nil
}

func (db *XcDB) doGetStrLen(key []byte) (uint32, error) {
	db.Mu.RLock()
	defer db.Mu.RUnlock()
	lsmMap := *db.Lsm
	tree := lsmMap[model.String]
	datainfo, err := tree.Get(key)
	if err != nil {
		return 0, err
	}
	offset := datainfo.Offset
	fileName := datainfo.FileName
	size := datainfo.Size

	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if isExpired(data) {
		err = errors.New("Data has expired")
		return 0, err
	}
	err = db.reSet(data)
	if err != nil {
		return 0, err
	}

	return datainfo.ValueSize, nil
}
func (db *XcDB) Append(key, value []byte) error {
	err := db.doAppend(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (db *XcDB) doAppend(key, value []byte) error {
	db.Mu.RLock()
	defer db.Mu.RUnlock()
	lsmMap := *db.Lsm
	tree := lsmMap[model.String]
	datainfo, err := tree.Get(key)
	if err != nil {
		return err
	}
	offset := datainfo.Offset
	fileName := datainfo.FileName
	size := datainfo.Size

	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)

	if err != nil {
		fmt.Println(err)
		return err
	}
	if isExpired(data) {
		err = errors.New("Data has expired")
		return err
	}
	data.DataMeta.Value = append(data.DataMeta.Value, value...)
	err = db.reSet(data)
	if err != nil {
		return err
	}
	return nil
}
