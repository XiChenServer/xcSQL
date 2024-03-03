package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/logs"
	"encoding/json"
	"time"
)

// 对于hash进行的建立
func (db *XcDB) Hset(key []byte, value map[string]string, ttl ...uint64) error {
	err := db.doHset(key, value, ttl...)
	return err
}
func (db *XcDB) doHset(key []byte, value map[string]string, ttl ...uint64) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	var timeSlice []time.Duration
	for _, t := range ttl {
		timeSlice = append(timeSlice, time.Duration(t)*time.Second)
	}

	lsmMap := *db.Lsm
	Hash := lsmMap[model.XCDB_Hash]

	values, err := mapToBytes(value)
	if err != nil {
		return err
	}

	e := NewKeyValueEntry(key, values, model.XCDB_Hash, model.XCDB_HSet, timeSlice...)
	stroeLocal, err := db.StorageManager.StoreData(e)
	if err != nil {
		logs.SugarLogger.Error("string set fail:", err)
		return err
	}

	datainfo := &lsm.DataInfo{
		DataMeta:        *e.DataMeta,
		StorageLocation: stroeLocal,
	}

	err = Hash.Insert(key, datainfo)
	return nil
}

// 将 map[string]string 转换为 []byte
func mapToBytes(m map[string]string) ([]byte, error) {
	// 将 map 转换为 JSON 格式的字符串
	jsonString, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	// 返回 JSON 字符串的字节表示
	return jsonString, nil
}

// 将 []byte 转换为 map[string]string
func bytesToMap(bytes []byte) (map[string]string, error) {
	// 解析 JSON 格式的字符串
	var m map[string]string
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	// 返回 map[string]string 类型的结果
	return m, nil
}
