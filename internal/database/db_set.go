package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"bytes"
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

// SRem 从集合中移除一个或多个成员
func (db *XcDB) SRem(key []byte, members [][]byte) error {
	err := db.doSRem(key, members)
	if err != nil {
		return err
	}
	return nil
}

func (db *XcDB) doSRem(key []byte, members [][]byte) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

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

	// 从集合中移除指定成员
	removedCount := 0
	for _, member := range members {
		for i, existingMember := range existingMembers {
			if bytes.Equal(existingMember, member) {
				// 从现有成员列表中删除匹配的成员
				existingMembers = append(existingMembers[:i], existingMembers[i+1:]...)
				removedCount++
				break
			}
		}
	}

	// 如果没有任何成员被移除，则直接返回
	if removedCount == 0 {
		return nil
	}

	// 将更新后的成员列表存储到数据库中
	dataValue := StoreSetMembers(existingMembers)
	e := NewKeyValueEntry(key, dataValue, model.XCDB_Set, model.XCDB_SetSREM)
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

// SMembers 获取集合中的所有成员
func (db *XcDB) SMembers(key []byte) ([][]byte, error) {
	return db.doSMembers(key)
}

func (db *XcDB) doSMembers(key []byte) ([][]byte, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

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
			return nil, err
		}

		if isExpired(data) {
			return nil, errors.New("Data has expired")
		}

		existingMembers, err = RetrieveSetMembers(data.Value)
		if err != nil {
			return nil, err
		}
	}

	return existingMembers, nil
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

// SIsMember 判断一个成员是否是集合的成员
func (db *XcDB) SIsMember(key []byte, member []byte) (bool, error) {
	return db.doSIsMember(key, member)
}

func (db *XcDB) doSIsMember(key []byte, member []byte) (bool, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	set := lsmMap[model.XCDB_Set]

	// 获取集合中的已有成员
	datainfo, _ := set.Get(key)
	if datainfo != nil {
		offset := datainfo.Offset
		fileName := datainfo.FileName
		size := datainfo.Size

		data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
		if err != nil {
			return false, err
		}

		if isExpired(data) {
			return false, errors.New("Data has expired")
		}

		existingMembers, err := RetrieveSetMembers(data.Value)
		if err != nil {
			return false, err
		}

		// 遍历集合中的成员，判断是否存在指定的成员
		for _, m := range existingMembers {
			if bytes.Equal(m, member) {
				return true, nil
			}
		}
	}

	return false, nil
}

// SCard 获取集合中成员的数量
func (db *XcDB) SCard(key []byte) (int, error) {
	return db.doSCard(key)
}

func (db *XcDB) doSCard(key []byte) (int, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	set := lsmMap[model.XCDB_Set]

	// 获取集合中的已有成员
	datainfo, _ := set.Get(key)
	if datainfo != nil {
		offset := datainfo.Offset
		fileName := datainfo.FileName
		size := datainfo.Size

		data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
		if err != nil {
			return 0, err
		}

		if isExpired(data) {
			return 0, errors.New("Data has expired")
		}

		existingMembers, err := RetrieveSetMembers(data.Value)
		if err != nil {
			return 0, err
		}

		// 返回集合中成员的数量
		return len(existingMembers), nil
	}

	// 如果集合不存在，返回数量为0
	return 0, nil
}
