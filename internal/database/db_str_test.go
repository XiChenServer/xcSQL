package database

import (
	"SQL/internal/storage"
	"SQL/logs"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// generateRandomData 生成指定长度的随机字节切片
func generateRandomData(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(rand.Intn(256)) // 生成0到255之间的随机字节
	}
	return data
}

// 简单的测试数据可以存入
func TestDB_Set(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()
	dataFilePath := "../../data/testdata/lsm_tree/test1.txt"

	// 加载模拟的数据文件到 LSM 树中
	err := db.lsm.LoadDataFromFile(dataFilePath)
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	key := []byte(generateRandomKey())
	value := []byte(generateRandomKey())
	err = db.Set(key, value)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Insert ok")
	fmt.Println(string(key))
	defer db.lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	defer db.lsm.SaveActiveToDiskOnExit()
	defer storage.SaveStorageManager(db.storageManager, "../../data/testdata/lsm_tree/config.txt")
}

// 简单的测试数据可以通过解压获取到
func TestDB_GetEasy(t *testing.T) {

	fileName := "data_0.gz" // 你的存储位置文件名
	offset := int64(222)    // 偏移量
	size := int64(227)      // 数据大小
	storage := &storage.StorageManager{}
	// 解压数据
	decompressedData, err := storage.DecompressAndFillData("../../data/testdata/string_test/"+fileName, offset, size)
	if err != nil {
		t.Fatalf("failed to decompress data: %v", err)
	}
	// 打印解压后的数据
	fmt.Println("Decompressed Data:", string(decompressedData.DataMeta.Key))
}

// 简单的测试数据可以存入
func TestDB_Get(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()

	// 加载模拟的数据文件到 LSM 树中
	err := db.lsm.LoadDataFromFile(string(db.lsm.LsmPath))
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	key := []byte("FtZlwNYtuo")
	data, err := db.Get(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Get ok:", string(data.DataMeta.Key), string(data.DataMeta.Value))
	defer db.lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	defer db.lsm.SaveActiveToDiskOnExit()
	defer storage.SaveStorageManager(db.storageManager, "../../data/testdata/lsm_tree/config.txt")
}

// generateRandomKey 生成随机键值
func generateRandomKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLen := 10
	b := make([]byte, keyLen)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
