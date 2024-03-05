package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/logs"
	"encoding/json"
	"errors"
	"time"
)

// 对于hash进行的建立
func (db *XcDB) HSet(key []byte, value map[string]string, ttl ...uint64) error {
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

// 对于hash进行的单个字段的获取
func (db *XcDB) HGet(key []byte, Field string) ([]byte, error) {
	data, err := db.doHGet(key, Field)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (db *XcDB) doHGet(key []byte, Field string) ([]byte, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	Hash := lsmMap[model.XCDB_Hash]
	datainfo, err := Hash.Get(key)
	if err != nil {
		return nil, err
	}
	offset := datainfo.Offset
	fileName := datainfo.FileName
	size := datainfo.Size
	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
	if err != nil {
		return nil, err
	}
	if isExpired(data) {
		err = errors.New("Data has expired")
		return nil, err
	}
	values, err := bytesToMap(data.Value)
	if err != nil {
		return nil, err
	}

	return []byte(values[Field]), nil
}

// 对于hash进行的所有字段的获取
func (db *XcDB) HGETALL(key []byte) (map[string]string, error) {
	data, err := db.doHGETALL(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (db *XcDB) doHGETALL(key []byte) (map[string]string, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	Hash := lsmMap[model.XCDB_Hash]
	datainfo, err := Hash.Get(key)
	if err != nil {
		return nil, err
	}
	offset := datainfo.Offset
	fileName := datainfo.FileName
	size := datainfo.Size
	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
	if err != nil {
		return nil, err
	}
	if isExpired(data) {
		err = errors.New("Data has expired")
		return nil, err
	}
	values, err := bytesToMap(data.Value)
	if err != nil {
		return nil, err
	}

	return values, nil
}

// HDel 用于删除哈希表中的一个或多个指定字段
func (db *XcDB) HDel(key []byte, fields ...string) error {
	err := db.doHDel(key, fields...)
	return err
}

func (db *XcDB) doHDel(key []byte, fields ...string) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	Hash := lsmMap[model.XCDB_Hash]

	// 获取哈希表中指定字段的数据信息
	dataInfo, err := Hash.Get(key)
	if err != nil {
		return err
	}

	// 从存储中读取数据
	offset := dataInfo.Offset
	fileName := dataInfo.FileName
	size := dataInfo.Size
	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
	if err != nil {
		return err
	}

	// 检查数据是否过期
	if isExpired(data) {
		err = errors.New("Data has expired")
		return err
	}

	// 将数据解析为map
	values, err := bytesToMap(data.Value)
	if err != nil {
		return err
	}

	// 删除指定字段
	for _, field := range fields {
		delete(values, field)
	}

	// 将更新后的数据重新序列化
	updatedValue, err := mapToBytes(values)
	if err != nil {
		return err
	}
	expiration := []time.Duration{dataInfo.TTL}

	// 创建新的数据条目
	e := NewKeyValueEntry(key, updatedValue, model.XCDB_Hash, model.XCDB_HSet, expiration...)

	// 更新存储中的数据
	storeLocal, err := db.StorageManager.StoreData(e)
	if err != nil {
		logs.SugarLogger.Error("string set fail:", err)
		return err
	}

	// 更新哈希表中的数据信息
	updatedDataInfo := &lsm.DataInfo{
		DataMeta:        *e.DataMeta,
		StorageLocation: storeLocal,
	}

	err = Hash.Insert(key, updatedDataInfo)
	if err != nil {
		return err
	}

	return nil
}

// HExists 用于检查哈希表中给定字段是否存在
func (db *XcDB) HExists(key []byte, field string) (bool, error) {
	exists, err := db.doHExists(key, field)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *XcDB) doHExists(key []byte, field string) (bool, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	Hash := lsmMap[model.XCDB_Hash]

	// 获取哈希表中指定字段的数据信息
	dataInfo, err := Hash.Get(key)
	if err != nil {
		return false, err
	}

	// 从存储中读取数据
	offset := dataInfo.Offset
	fileName := dataInfo.FileName
	size := dataInfo.Size
	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
	if err != nil {
		return false, err
	}

	// 检查数据是否过期
	if isExpired(data) {
		err = errors.New("Data has expired")
		return false, err
	}

	// 将数据解析为map
	values, err := bytesToMap(data.Value)
	if err != nil {
		return false, err
	}

	// 检查字段是否存在
	_, exists := values[field]

	return exists, nil
}

// HKeys 用于获取哈希表中的所有字段
func (db *XcDB) HKeys(key []byte) ([]string, error) {
	keys, err := db.doHKeys(key)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (db *XcDB) doHKeys(key []byte) ([]string, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	Hash := lsmMap[model.XCDB_Hash]

	// 获取哈希表中指定字段的数据信息
	dataInfo, err := Hash.Get(key)
	if err != nil {
		return nil, err
	}

	// 从存储中读取数据
	offset := dataInfo.Offset
	fileName := dataInfo.FileName
	size := dataInfo.Size
	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
	if err != nil {
		return nil, err
	}

	// 检查数据是否过期
	if isExpired(data) {
		err = errors.New("Data has expired")
		return nil, err
	}

	// 将数据解析为map
	values, err := bytesToMap(data.Value)
	if err != nil {
		return nil, err
	}

	// 提取所有字段名
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}

	return keys, nil
}

// HVals 用于获取哈希表中的所有值
func (db *XcDB) HVals(key []byte) ([]string, error) {
	vals, err := db.doHVals(key)
	if err != nil {
		return nil, err
	}
	return vals, nil
}

func (db *XcDB) doHVals(key []byte) ([]string, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	Hash := lsmMap[model.XCDB_Hash]

	// 获取哈希表中指定字段的数据信息
	dataInfo, err := Hash.Get(key)
	if err != nil {
		return nil, err
	}

	// 从存储中读取数据
	offset := dataInfo.Offset
	fileName := dataInfo.FileName
	size := dataInfo.Size
	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
	if err != nil {
		return nil, err
	}

	// 检查数据是否过期
	if isExpired(data) {
		err = errors.New("Data has expired")
		return nil, err
	}

	// 将数据解析为map
	values, err := bytesToMap(data.Value)
	if err != nil {
		return nil, err
	}

	// 提取所有值
	var vals []string
	for _, val := range values {
		vals = append(vals, val)
	}

	return vals, nil
}

// HLen 用于获取哈希表中字段的数量
func (db *XcDB) HLen(key []byte) (int, error) {
	length, err := db.doHLen(key)
	if err != nil {
		return 0, err
	}
	return length, nil
}

func (db *XcDB) doHLen(key []byte) (int, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	Hash := lsmMap[model.XCDB_Hash]

	// 获取哈希表中指定字段的数据信息
	dataInfo, err := Hash.Get(key)
	if err != nil {
		return 0, err
	}

	// 从存储中读取数据
	offset := dataInfo.Offset
	fileName := dataInfo.FileName
	size := dataInfo.Size
	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
	if err != nil {
		return 0, err
	}

	// 检查数据是否过期
	if isExpired(data) {
		err = errors.New("Data has expired")
		return 0, err
	}

	// 将数据解析为map
	values, err := bytesToMap(data.Value)
	if err != nil {
		return 0, err
	}

	// 计算字段的数量
	length := len(values)

	return length, nil
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
