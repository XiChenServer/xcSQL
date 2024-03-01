package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/logs"
	"bytes"
	"encoding/binary"
	"time"
)

// 编写对于列表的操作方法
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
	changeValue := StoreListValueWithDataType(values)
	e := NewKeyValueEntry(key, changeValue, model.List, model.ListLPUSH, timeSlice...)
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
	tree := lsmMap[model.List]
	tree.Insert(key, datainfo)
	return nil
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
