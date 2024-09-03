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

//
//import (
//	"SQL/internal/model"
//	"SQL/internal/storage"
//	"SQL/logs"
//	"fmt"
//	"math/rand"
//	"reflect"
//	"testing"
//	"time"
//)

func Test_HSET(t *testing.T) {
	logs.InitLogger("")

	dbPool := NewConnectionPool(2, 10, 30*time.Minute)

	driver, _ := dbPool.GetConnection("123")
	fmt.Println(123)
	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
	lsmMap := *driver.DB.Lsm
	lsmType := lsmMap[model.XCDB_Hash]
	// 加载模拟的数据文件到 LSM 树中
	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
	//if err != nil {
	//	t.Fatalf("Error loading data from disk: %v", err)
	//}
	//key := []byte("UDVGKnSAsp")
	//key := []byte(generateRandomKey())
	//value := generateRandomMap(4, 5)

	key := []byte("people")
	myMap := map[string]string{
		"name":    "John",
		"age":     "30",
		"address": "123 Main St",
	}

	err := driver.DB.HSet(key, myMap)
	if err != nil {
		return
	}

	fmt.Println("SET ok")

	p, _ := driver.DB.HGet([]byte("people"), "name")
	fmt.Println(string(p))
	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(driver.DB.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}

// generateRandomString 生成指定长度的随机字符串
func generateRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// generateRandomMap 生成随机的键值对映射
func generateRandomMap(size, valueLength int) map[string]string {
	result := make(map[string]string, size)
	for i := 0; i < size; i++ {
		key := generateRandomString(8)
		value := generateRandomString(valueLength)
		result[key] = value
	}
	return result
}

func Test_HSET12(t *testing.T) {
	logs.InitLogger("")

	dbPool := NewConnectionPool(2, 10, 30*time.Minute)

	driver, err := dbPool.GetConnection("123")
	if err != nil {
		t.Fatalf("failed to get connection: %v", err)
	}
	defer driver.DB.Close()

	lsmMap := *driver.DB.Lsm
	lsmType := lsmMap[model.XCDB_Hash]

	var wg sync.WaitGroup
	concurrencyLevel := 10 // 设置并发级别

	// 启动多个 goroutine 进行并发写入
	for i := 0; i < concurrencyLevel; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := []byte(generateRandomString(10))
			valueMap := generateRandomMap(5, 10) // 生成一个包含5个键值对的随机 map
			err := driver.DB.HSet(key, valueMap)
			if err != nil {
				t.Errorf("Goroutine %d: failed to set value: %v", i, err)
			} else {
				fmt.Printf("Goroutine %d: SET ok for key: %s\n", i, key)
			}
		}(i)
	}

	// 启动多个 goroutine 进行并发读取
	for i := 0; i < concurrencyLevel; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := []byte(generateRandomString(10))
			field := generateRandomString(5)
			p, err := driver.DB.HGet(key, field)
			if err != nil {
				t.Errorf("Goroutine %d: failed to get value: %v", i, err)
			} else {
				fmt.Printf("Goroutine %d: GET value: %s for key: %s, field: %s\n", i, p, key, field)
			}
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 保存数据并关闭数据库
	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(driver.DB.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}

//func Test_HGET(t *testing.T) {
//	logs.InitLogger("")
//	db, _ := NewXcDB("123", NewConnectionPool(2, 10, 30*time.Minute))
//	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
//	lsmMap := *db.Lsm
//	lsmType := lsmMap[model.XCDB_Hash]
//	// 加载模拟的数据文件到 LSM 树中
//	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
//	//if err != nil {
//	//	t.Fatalf("Error loading data from disk: %v", err)
//	//}
//	//key := []byte("UDVGKnSAsp")
//	//key := []byte(generateRandomKey())
//	//value := generateRandomMap(4, 5)
//	key := []byte("people")
//	field := "1232"
//	data, err := db.HGet(key, field)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println("HGET ok")
//	fmt.Println(string(data))
//	lsmType.SaveActiveToDiskOnExit()
//	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
//	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
//}
//
//func Test_HGETALL(t *testing.T) {
//	logs.InitLogger("")
//	db, _ := NewXcDB("123", NewConnectionPool(2, 10, 30*time.Minute))
//	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
//	lsmMap := *db.Lsm
//	lsmType := lsmMap[model.XCDB_Hash]
//	// 加载模拟的数据文件到 LSM 树中
//	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
//	//if err != nil {
//	//	t.Fatalf("Error loading data from disk: %v", err)
//	//}
//	//key := []byte("UDVGKnSAsp")
//	//key := []byte(generateRandomKey())
//	//value := generateRandomMap(4, 5)
//	key := []byte("people")
//	//field := "1232"
//	data, err := db.HGETALL(key)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println("HGET ok")
//	fmt.Println(data)
//	lsmType.SaveActiveToDiskOnExit()
//	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
//	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
//}
//
//func Test_MapToBytesAndBytesToMap(t *testing.T) {
//	// 定义一个 map[string]string
//	value := generateRandomMap(4, 5)
//	// 将 map[string]string 转换为 []byte
//	bytes, err := mapToBytes(value)
//	if err != nil {
//		t.Errorf("mapToBytes returned an error: %v", err)
//	}
//
//	// 将 []byte 转换为 map[string]string
//	m, err := bytesToMap(bytes)
//	if err != nil {
//		t.Errorf("bytesToMap returned an error: %v", err)
//	}
//
//	// 检查原始 map 和解析出来的 map 是否相等
//	if !reflect.DeepEqual(value, m) {
//		t.Errorf("mapToBytes and bytesToMap result mismatch: %v != %v", value, m)
//	}
//}
//
//func Test_HDEL(t *testing.T) {
//	logs.InitLogger("")
//	db, _ := NewXcDB("123", NewConnectionPool(2, 10, 30*time.Minute))
//	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
//	lsmMap := *db.Lsm
//	lsmType := lsmMap[model.XCDB_Hash]
//	// 加载模拟的数据文件到 LSM 树中
//	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
//	//if err != nil {
//	//	t.Fatalf("Error loading data from disk: %v", err)
//	//}
//	//key := []byte("UDVGKnSAsp")
//	//key := []byte(generateRandomKey())
//	//value := generateRandomMap(4, 5)
//
//	key := []byte("people")
//	//myMap := map[string]string{
//	//	"name":    "John",
//	//	"age":     "30",
//	//	"address": "123 Main St",
//	//}
//
//	err := db.HDel(key, "name")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println("HSET ok")
//	lsmType.SaveActiveToDiskOnExit()
//	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
//	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
//}
//
//func Test_HExists(t *testing.T) {
//	logs.InitLogger("")
//	db, _ := NewXcDB("123", NewConnectionPool(2, 10, 30*time.Minute))
//	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
//	lsmMap := *db.Lsm
//	lsmType := lsmMap[model.XCDB_Hash]
//	// 加载模拟的数据文件到 LSM 树中
//	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
//	//if err != nil {
//	//	t.Fatalf("Error loading data from disk: %v", err)
//	//}
//	//key := []byte("UDVGKnSAsp")
//	//key := []byte(generateRandomKey())
//	//value := generateRandomMap(4, 5)
//
//	key := []byte("people")
//	//myMap := map[string]string{
//	//	"name":    "John",
//	//	"age":     "30",
//	//	"address": "123 Main St",
//	//}
//
//	flag, err := db.HExists(key, "age")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(flag)
//	fmt.Println("HSET ok")
//	lsmType.SaveActiveToDiskOnExit()
//	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
//	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
//}
//
//func Test_HKeys(t *testing.T) {
//	logs.InitLogger("")
//	db, _ := NewXcDB("123", NewConnectionPool(2, 10, 30*time.Minute))
//	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
//	lsmMap := *db.Lsm
//	lsmType := lsmMap[model.XCDB_Hash]
//	// 加载模拟的数据文件到 LSM 树中
//	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
//	//if err != nil {
//	//	t.Fatalf("Error loading data from disk: %v", err)
//	//}
//	//key := []byte("UDVGKnSAsp")
//	//key := []byte(generateRandomKey())
//	//value := generateRandomMap(4, 5)
//
//	key := []byte("people")
//	//myMap := map[string]string{
//	//	"name":    "John",
//	//	"age":     "30",
//	//	"address": "123 Main St",
//	//}
//
//	flag, err := db.HKeys(key)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(flag)
//	fmt.Println("HSET ok")
//	lsmType.SaveActiveToDiskOnExit()
//	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
//	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
//}
//
//func Test_HVals(t *testing.T) {
//	logs.InitLogger("")
//	db, _ := NewXcDB("123", NewConnectionPool(2, 10, 30*time.Minute))
//	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
//	lsmMap := *db.Lsm
//	lsmType := lsmMap[model.XCDB_Hash]
//	// 加载模拟的数据文件到 LSM 树中
//	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
//	//if err != nil {
//	//	t.Fatalf("Error loading data from disk: %v", err)
//	//}
//	//key := []byte("UDVGKnSAsp")
//	//key := []byte(generateRandomKey())
//	//value := generateRandomMap(4, 5)
//
//	key := []byte("people")
//	//myMap := map[string]string{
//	//	"name":    "John",
//	//	"age":     "30",
//	//	"address": "123 Main St",
//	//}
//
//	flag, err := db.HVals(key)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(flag)
//	fmt.Println("HSET ok")
//	lsmType.SaveActiveToDiskOnExit()
//	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
//	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
//}
//
//func Test_HLen(t *testing.T) {
//	logs.InitLogger("")
//	db, _ := NewXcDB("123", NewConnectionPool(2, 10, 30*time.Minute))
//	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
//	lsmMap := *db.Lsm
//	lsmType := lsmMap[model.XCDB_Hash]
//	// 加载模拟的数据文件到 LSM 树中
//	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
//	//if err != nil {
//	//	t.Fatalf("Error loading data from disk: %v", err)
//	//}
//	//key := []byte("UDVGKnSAsp")
//	//key := []byte(generateRandomKey())
//	//value := generateRandomMap(4, 5)
//
//	key := []byte("people")
//	//myMap := map[string]string{
//	//	"name":    "John",
//	//	"age":     "30",
//	//	"address": "123 Main St",
//	//}
//
//	flag, err := db.HLen(key)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(flag)
//	fmt.Println("HSET ok")
//	lsmType.SaveActiveToDiskOnExit()
//	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
//	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
//}
//
//// 随机生成指定长度的字符串
//func randomString(length int) string {
//	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
//	b := make([]byte, length)
//	for i := range b {
//		b[i] = charset[rand.Intn(len(charset))]
//	}
//	return string(b)
//}
//
//// 生成随机的 map[string]string 数据
//func generateRandomMap(length, size int) map[string]string {
//	randomMap := make(map[string]string)
//	for i := 0; i < size; i++ {
//		key := randomString(length)
//		value := randomString(length)
//		randomMap[key] = value
//	}
//	return randomMap
//}
