package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/logs"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"
)

// 编写对于列表的操作方法

// 注意这里的操作是对于list的尾部进行插入的操作
func (db *XcDB) RPUSH(key []byte, values [][]byte, ttl ...uint64) error {
	err := db.doRPUSH(key, values, ttl...)
	if err != nil {
		return err
	}
	return nil
}

func (db *XcDB) doRPUSH(key []byte, values [][]byte, ttl ...uint64) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	var timeSlice []time.Duration
	for _, t := range ttl {
		timeSlice = append(timeSlice, time.Duration(t)*time.Second)
	}

	lsmMap := *db.Lsm
	list := lsmMap[model.XCDB_List]
	datainfo, _ := list.Get(key)
	if datainfo != nil {
		offset := datainfo.Offset
		fileName := datainfo.FileName
		size := datainfo.Size

		data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if isExpired(data) {
			return errors.New("Data has expired")
		}

		oldValue, err := RetrieveListValueWithDataType(data.Value)

		//
		// 将新值插入到旧值的后面
		for _, item := range values {
			oldValue = append(oldValue, item)
		}

		changeValue := StoreListValueWithDataType(oldValue)

		//fmt.Println(changeValue)
		//data.Value = append(data.Value, changeValue...)

		e := NewKeyValueEntry(key, changeValue, model.XCDB_List, model.XCDB_ListLPUSH, timeSlice...)
		stroeLocal, err := db.StorageManager.StoreData(e)
		if err != nil {
			logs.SugarLogger.Error("string set fail:", err)
			return err
		}

		datainfo = &lsm.DataInfo{
			DataMeta:        *e.DataMeta,
			StorageLocation: stroeLocal,
		}

		err = list.Insert(key, datainfo)
		return err
	}

	changeValue := StoreListValueWithDataType(values)
	e := NewKeyValueEntry(key, changeValue, model.XCDB_List, model.XCDB_ListLPUSH, timeSlice...)
	//fmt.Println(e.Value)
	stroeLocal, err := db.StorageManager.StoreData(e)
	if err != nil {
		logs.SugarLogger.Error("string set fail:", err)
		return err
	}

	datainfo = &lsm.DataInfo{
		DataMeta:        *e.DataMeta,
		StorageLocation: stroeLocal,
	}

	err = list.Insert(key, datainfo)
	if err != nil {
		return err
	}
	return nil
}

// 注意这里的操作是对于list的头部进行插入的操作
func (db *XcDB) LPUSH(key []byte, values [][]byte, ttl ...uint64) error {
	err := db.doLPUSH(key, values, ttl...)
	if err != nil {
		return err
	}
	return nil
}

func (db *XcDB) doLPUSH(key []byte, values [][]byte, ttl ...uint64) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	var timeSlice []time.Duration
	for _, t := range ttl {
		timeSlice = append(timeSlice, time.Duration(t)*time.Second)
	}

	lsmMap := *db.Lsm
	list := lsmMap[model.XCDB_List]
	datainfo, _ := list.Get(key)

	// 创建新的列表并将值插入其中
	if datainfo == nil {
		// 对于新创建的列表，需要将插入的值逆序存储
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
		changeValue := StoreListValueWithDataType(values)
		e := NewKeyValueEntry(key, changeValue, model.XCDB_List, model.XCDB_ListLPUSH, timeSlice...)
		stroeLocal, err := db.StorageManager.StoreData(e)
		if err != nil {
			logs.SugarLogger.Error("string set fail:", err)
			return err
		}

		datainfo = &lsm.DataInfo{
			DataMeta:        *e.DataMeta,
			StorageLocation: stroeLocal,
		}

		err = list.Insert(key, datainfo)
		return err
	}

	// 如果列表存在，则将新值插入到旧值的前面
	offset := datainfo.Offset
	fileName := datainfo.FileName
	size := datainfo.Size

	data, err := db.StorageManager.DecompressAndFillData(string(fileName), offset, size)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if isExpired(data) {
		return errors.New("Data has expired")
	}

	oldValue, err := RetrieveListValueWithDataType(data.Value)
	if err != nil {
		return err
	}
	// 对于后续追加的值，同样需要逆序存储
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	// 将新值插入到旧值的前面
	// 由于是左侧插入，所以要先将旧值插入到新值前面，然后再将新值插入到新值的前面
	values = append(values, oldValue...)

	// 存储修改后的列表值
	changeValue := StoreListValueWithDataType(values)
	e := NewKeyValueEntry(key, changeValue, model.XCDB_List, model.XCDB_ListLPUSH, timeSlice...)
	stroeLocal, err := db.StorageManager.StoreData(e)
	if err != nil {
		logs.SugarLogger.Error("string set fail:", err)
		return err
	}

	datainfo = &lsm.DataInfo{
		DataMeta:        *e.DataMeta,
		StorageLocation: stroeLocal,
	}

	err = list.Insert(key, datainfo)
	return err
}

// 对于list进行查找的操作，
func (db *XcDB) LRANGE(key []byte, left, right int) ([][]byte, error) {
	if (left == 0 && right == 0) || math.Abs(float64(left)) > math.Abs(float64(right)) {
		return nil, errors.New("both left and right values are required")
	}
	data, err := db.doLRANGE(key, left, right)
	return data, err
}

func (db *XcDB) doLRANGE(key []byte, left, right int) ([][]byte, error) {
	db.Mu.RLock()
	defer db.Mu.RUnlock()
	lsmMap := *db.Lsm
	tree := lsmMap[model.XCDB_List]
	datainfo, err := tree.Get(key)
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
	value, err := getValueByRange(data.Value, left, right)
	err = db.reSet(data)
	if err != nil {
		return nil, err
	}

	if data == nil {
		err := errors.New("No data found")
		return nil, err
	}
	return value, nil
}

// 对于list进行查找的操作，
func (db *XcDB) LINDEX(key []byte, index int) ([]byte, error) {
	if index < 0 {
		return nil, errors.New("both left and right values are required")
	}
	data, err := db.doLINDEX(key, index)
	return data, err
}

func (db *XcDB) doLINDEX(key []byte, index int) ([]byte, error) {
	db.Mu.RLock()
	defer db.Mu.RUnlock()
	lsmMap := *db.Lsm
	tree := lsmMap[model.XCDB_List]
	datainfo, err := tree.Get(key)
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
	value, err := getValueAtIndex(data.Value, index)
	err = db.reSet(data)
	if err != nil {
		return nil, err
	}

	if data == nil {
		err := errors.New("No data found")
		return nil, err
	}
	return value, nil
}

// 根据下标获取值
func getValueAtIndex(data []byte, index int) ([]byte, error) {
	values, err := RetrieveListValueWithDataType(data)
	if err != nil {
		return nil, err
	}
	for _, v := range values {
		fmt.Println(string(v))
	}
	if index < 0 || index >= len(values) {
		return nil, errors.New("index out of range")
	}

	return values[index], nil
}

// 根据左右的范围获取值
func getValueByRange(data []byte, left, right int) ([][]byte, error) {

	value, err := RetrieveListValueWithDataType(data)
	if err != nil {
		return nil, err
	}

	// 如果 left 为 0，right 为 -1，返回所有数据
	if left == 0 && right == -1 {
		return value, nil
	}

	// 如果 left 或 right 超出范围，则返回空
	if left >= len(value) || right >= len(value) {
		return nil, nil
	}

	// 根据 left 和 right 范围取值
	return value[left : right+1], nil
}

// 存储列表类型的值，使用标记区分数据类型
func StoreListValueWithDataType(elements [][]byte) []byte {
	var dataBuffer bytes.Buffer
	// 写入数据类型标记
	binary.Write(&dataBuffer, binary.BigEndian, uint16(1)) // 使用1表示列表类型
	// 写入列表元素
	for _, element := range elements {
		// 写入元素长度
		binary.Write(&dataBuffer, binary.BigEndian, uint32(len(element)))
		// 写入元素值
		dataBuffer.Write(element)
	}
	return dataBuffer.Bytes()
}

// 从存储中取出列表类型的值，并解析数据类型
func RetrieveListValueWithDataType(data []byte) ([][]byte, error) {
	var elements [][]byte
	// 读取数据类型标记
	buffer := bytes.NewReader(data)
	var valueType uint16
	err := binary.Read(buffer, binary.BigEndian, &valueType)
	if err != nil {
		return nil, err
	}
	// 解析列表元素
	for buffer.Len() > 0 {
		// 读取元素长度
		var elementLength uint32
		err := binary.Read(buffer, binary.BigEndian, &elementLength)
		if err != nil {
			return nil, err
		}
		// 读取元素值
		element := make([]byte, elementLength)
		_, err = buffer.Read(element)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}
	return elements, nil
}

// LPOP 从列表中移除并返回头部元素
func (db *XcDB) LPOP(key []byte) ([]byte, error) {
	data, err := db.doLPOP(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (db *XcDB) doLPOP(key []byte) ([]byte, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	list := lsmMap[model.XCDB_List]
	datainfo, err := list.Get(key)
	if err != nil {
		return nil, err
	}
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

		// 解析列表值
		oldValue, err := RetrieveListValueWithDataType(data.Value)
		if err != nil {
			return nil, err
		}

		// 移除并返回头部元素
		var head []byte
		if len(oldValue) > 0 {
			head = oldValue[0]
			oldValue = oldValue[1:]
		} else {
			// 如果列表为空，则返回空
			return nil, errors.New("List is empty")
		}

		// 更新数据并存储
		changeValue := StoreListValueWithDataType(oldValue)
		e := NewKeyValueEntry(key, changeValue, model.XCDB_List, model.XCDB_ListLPOP)
		storeLocal, err := db.StorageManager.StoreData(e)
		if err != nil {
			return nil, err
		}

		datainfo = &lsm.DataInfo{
			DataMeta:        *e.DataMeta,
			StorageLocation: storeLocal,
		}

		err = list.Insert(key, datainfo)
		if err != nil {
			return nil, err
		}

		return head, nil
	}
	return nil, errors.New("List not found")
}

// RPOP 从列表中移除并返回尾部元素
func (db *XcDB) RPOP(key []byte) ([]byte, error) {
	data, err := db.doRPOP(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// doRPOP 实际执行RPOP操作
func (db *XcDB) doRPOP(key []byte) ([]byte, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	lsmMap := *db.Lsm
	list := lsmMap[model.XCDB_List]
	datainfo, err := list.Get(key)
	if err != nil {
		return nil, err
	}
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

		// 解析列表值
		oldValue, err := RetrieveListValueWithDataType(data.Value)
		if err != nil {
			return nil, err
		}

		// 移除并返回尾部元素
		var tail []byte
		if len(oldValue) > 0 {
			lastIndex := len(oldValue) - 1
			tail = oldValue[lastIndex]
			oldValue = oldValue[:lastIndex]
		} else {
			// 如果列表为空，则返回空
			return nil, errors.New("List is empty")
		}

		// 更新数据并存储
		changeValue := StoreListValueWithDataType(oldValue)
		e := NewKeyValueEntry(key, changeValue, model.XCDB_List, model.XCDB_ListRPOP)
		storeLocal, err := db.StorageManager.StoreData(e)
		if err != nil {
			return nil, err
		}

		datainfo = &lsm.DataInfo{
			DataMeta:        *e.DataMeta,
			StorageLocation: storeLocal,
		}

		err = list.Insert(key, datainfo)
		if err != nil {
			return nil, err
		}

		return tail, nil
	}
	return nil, errors.New("List not found")
}
