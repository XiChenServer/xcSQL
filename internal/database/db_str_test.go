package database

import (
	"SQL/internal/model"
	"SQL/internal/storage"
	"SQL/logs"
	"fmt"
	"math/rand"
	"sync"
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
func TestDB_SetMore(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()
	dataFilePath := "../../data/testdata/lsm_tree/test1.txt"

	// 加载模拟的数据文件到 LSM 树中
	lsmMap := *db.Lsm
	tree := lsmMap[model.XCDB_String]
	err := tree.LoadDataFromFile(dataFilePath)
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}

	// 并发插入的总数量
	numInserts := 1000

	// 用于等待所有 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(numInserts)

	// 模拟并发插入操作
	for i := 0; i < numInserts; i++ {
		go func(i int) {
			defer wg.Done() // 表示当前 goroutine 已完成

			key := []byte(generateRandomKey())
			value := []byte(generateRandomKey())

			err = db.Set(key, value)

		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	tree.SaveActiveToDiskOnExit()
	tree.PrintDiskDataToFile(string(tree.LsmPath))
	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}

// 简单的测试数据可以存入
func TestDB_Set(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()
	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
	lsmMap := *db.Lsm
	lsmType := lsmMap[model.XCDB_String]
	// 加载模拟的数据文件到 LSM 树中
	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
	//if err != nil {
	//	t.Fatalf("Error loading data from disk: %v", err)
	//}
	key := []byte(generateRandomKey())
	value := []byte(generateRandomKey())
	err := db.Set(key, value)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Insert ok")
	fmt.Println(string(key), string(value))

	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
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
	lsmMap := *db.Lsm
	tree := lsmMap[model.XCDB_String]
	// 加载模拟的数据文件到 LSM 树中
	err := tree.LoadDataFromFile(string(tree.LsmPath))
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	key := []byte("LtAkGhMNFf")
	data, err := db.Get(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Get ok:", string(data.DataMeta.Key), string(data.Value))
	defer tree.PrintDiskDataToFile(string(tree.LsmPath))
	defer tree.SaveActiveToDiskOnExit()
	defer storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
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

// 简单的测试数据可以获取长度
func TestDB_Strlen(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()
	lsmMap := *db.Lsm
	tree := lsmMap[model.XCDB_String]
	// 加载模拟的数据文件到 LSM 树中
	err := tree.LoadDataFromFile(string(tree.LsmPath))
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	key := []byte("Lykig5qLNL")
	data, err := db.Strlen(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Get ok:", data)
	defer tree.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	defer tree.SaveActiveToDiskOnExit()
	defer storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}

// 进行追加操作
func TestDB_Append(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()
	lsmMap := *db.Lsm
	tree := lsmMap[model.XCDB_String]
	// 加载模拟的数据文件到 LSM 树中
	err := tree.LoadDataFromFile(string(tree.LsmPath))
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	key := []byte("NG9tAX3q3V")
	err = db.Append(key, []byte("dfdsfd"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Append success")
	defer tree.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	defer tree.SaveActiveToDiskOnExit()
	defer storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}
