package database

import (
	"SQL/internal/model"
	"SQL/internal/storage"
	"SQL/logs"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// 简单的测试数据可以存入
func TestDB_LPUSH(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()
	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
	lsmMap := *db.Lsm
	lsmType := lsmMap[model.List]
	// 加载模拟的数据文件到 LSM 树中
	err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	key := []byte(generateRandomKey())
	value := generateRandomByteSlices(4, 6, 5)
	err = db.LPUSH(key, value)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Insert ok")
	fmt.Println(string(key))
	fmt.Println(string(lsmType.LsmPath))

	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}

// 测试将[][]byte类型，压缩成[]byte类型
func Test_value_change(t *testing.T) {
	// 生成随机数据
	values := generateRandomByteSlices(10, 10, 5)

	// 存储数据
	value := StoreListValueWithDataType(values)
	fmt.Println((value))
	// 检索存储的数据
	next, err := RetrieveListValueWithDataType(value)
	if err != nil {
		t.Errorf("Error retrieving list value: %v", err)
	}

	// 比较原始数据和检索到的数据是否相同
	if !reflect.DeepEqual(values, next) {
		t.Errorf("Retrieved list value does not match original value.")
	}

}

// 随机生成 [][]byte 类型的数据
func generateRandomByteSlices(n, maxElementSize, maxSliceLength int) [][]byte {
	rand.Seed(time.Now().UnixNano())

	var result [][]byte

	for i := 0; i < n; i++ {
		// 随机生成切片长度
		sliceLength := rand.Intn(maxSliceLength) + 1
		// 随机生成切片元素
		var byteSlice []byte
		for j := 0; j < sliceLength; j++ {
			// 随机生成元素长度
			elementSize := rand.Intn(maxElementSize) + 1
			// 随机生成元素值
			element := make([]byte, elementSize)
			rand.Read(element)
			byteSlice = append(byteSlice, element...)
		}
		result = append(result, byteSlice)
	}

	return result
}