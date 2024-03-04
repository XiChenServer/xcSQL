package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"encoding/json"
	"errors"
	"time"
)

// SAdd 将一个或多个成员添加到集合中
func (db *XcDB) SAdd(key []byte, members [][]byte, ttl ...uint64) error {
	err := db.doSAdd(key, members, ttl...)
	if err != nil {
		return err
	}
	return nil
}

func (db *XcDB) doSAdd(key []byte, members [][]byte, ttl ...uint64) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	var timeSlice []time.Duration
	for _, t := range ttl {
		timeSlice = append(timeSlice, time.Duration(t)*time.Second)
	}

	lsmMap := *db.Lsm
	set := lsmMap[model.XCDB_Set]

	// 获取集合中的已有成员
	datainfo, _ := set.Get(key)
	var existingMembers [][]byte
	if datainfo != nil {
		offset := datainfo.Offset
		fileName := datainfo.FileName
		size := datainfo.Size

		data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
		if err != nil {
			return err
		}

		if isExpired(data) {
			return errors.New("Data has expired")
		}

		existingMembers, err = RetrieveSetMembers(data.Value)
		if err != nil {
			return err
		}
	}

	// 将新成员添加到已有成员中
	for _, member := range members {
		existingMembers = append(existingMembers, member)
	}

	// 移除重复的成员
	uniqueMembers := make(map[string]bool)
	for _, member := range existingMembers {
		uniqueMembers[string(member)] = true
	}

	var uniqueMembersSlice [][]byte
	for member := range uniqueMembers {
		uniqueMembersSlice = append(uniqueMembersSlice, []byte(member))
	}

	// 将更新后的成员列表存储到数据库中
	dataValue := StoreSetMembers(uniqueMembersSlice)
	e := NewKeyValueEntry(key, dataValue, model.XCDB_Set, model.XCDB_SetSADD, timeSlice...)
	storeLocal, err := db.StorageManager.StoreData(e)
	if err != nil {
		return err
	}

	datainfo = &lsm.DataInfo{
		DataMeta:        *e.DataMeta,
		StorageLocation: storeLocal,
	}

	// 更新集合中的数据信息
	err = set.Insert(key, datainfo)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveSetMembers 用于从字节数组中检索集合的成员
func RetrieveSetMembers(data []byte) ([][]byte, error) {
	var members [][]byte
	err := json.Unmarshal(data, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// StoreSetMembers 将集合成员编码为字节数组
func StoreSetMembers(members [][]byte) []byte {
	encodedMembers, _ := json.Marshal(members)
	return encodedMembers
}
