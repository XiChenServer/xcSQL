package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/logs"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// 编写对于列表的操作方法

// 注意这里的操作是对于list的头部进行插入的操作
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

		// 创建一个新的切片用于拼接
		var combined [][]byte

		// 遍历第一个切片并将其内容逐个添加到新切片中
		for _, item := range oldValue {
			combined = append(combined, item)
		}

		// 遍历第二个切片并将其内容逐个添加到新切片中
		for _, item := range values {
			combined = append(combined, item)
		}
		changeValue := StoreListValueWithDataType(combined)

		fmt.Println(changeValue)
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
	fmt.Println(e.Value)
	stroeLocal, err := db.StorageManager.StoreData(e)
	if err != nil {
		logs.SugarLogger.Error("string set fail:", err)
		return err
	}

	datainfo = &lsm.DataInfo{
		DataMeta:        *e.DataMeta,
		StorageLocation: stroeLocal,
	}

	list.Insert(key, datainfo)
	return nil
}

// 对于list进行查找的操作，
func (db *XcDB) LRANGE(key []byte, left, right int) ([][]byte, error) {
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

func getValueByRange(data []byte, left, rignt int) ([][]byte, error) {
	value, err := RetrieveListValueWithDataType(data)
	if err != nil {
		return nil, err
	}
	var string1 []string
	for _, b := range value {
		fmt.Println(string(b))
		string1 = append(string1, string(b))
	}

	//// 使用 strings.Join 将 []string 拼接成一个字符串并打印
	//result := fmt.Sprintf("[%s]", strings.Join(string1, ", "))
	//fmt.Println(result)
	return nil, nil
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
