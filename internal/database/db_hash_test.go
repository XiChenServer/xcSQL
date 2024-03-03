package database

import (
	"SQL/internal/model"
	"SQL/internal/storage"
	"SQL/logs"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func Test_HSET(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()
	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
	lsmMap := *db.Lsm
	lsmType := lsmMap[model.XCDB_Hash]
	// 加载模拟的数据文件到 LSM 树中
	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
	//if err != nil {
	//	t.Fatalf("Error loading data from disk: %v", err)
	//}
	//key := []byte("UDVGKnSAsp")
	key := []byte(generateRandomKey())
	value := generateRandomMap(4, 5)
	err := db.Hset(key, value)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("HSET ok")
	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}

func Test_MapToBytesAndBytesToMap(t *testing.T) {
	// 定义一个 map[string]string
	value := generateRandomMap(4, 5)
	// 将 map[string]string 转换为 []byte
	bytes, err := mapToBytes(value)
	if err != nil {
		t.Errorf("mapToBytes returned an error: %v", err)
	}

	// 将 []byte 转换为 map[string]string
	m, err := bytesToMap(bytes)
	if err != nil {
		t.Errorf("bytesToMap returned an error: %v", err)
	}

	// 检查原始 map 和解析出来的 map 是否相等
	if !reflect.DeepEqual(value, m) {
		t.Errorf("mapToBytes and bytesToMap result mismatch: %v != %v", value, m)
	}
}

// 随机生成指定长度的字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// 生成随机的 map[string]string 数据
func generateRandomMap(length, size int) map[string]string {
	randomMap := make(map[string]string)
	for i := 0; i < size; i++ {
		key := randomString(length)
		value := randomString(length)
		randomMap[key] = value
	}
	return randomMap
}
