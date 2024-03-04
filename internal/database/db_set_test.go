package database

import (
	"SQL/internal/model"
	"SQL/internal/storage"
	"SQL/logs"
	"fmt"
	"testing"
)

// 简单的测试数据可以存入
func TestDB_SADD(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()
	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
	lsmMap := *db.Lsm
	lsmType := lsmMap[model.XCDB_Set]
	//// 加载模拟的数据文件到 LSM 树中
	//err := lsmType.LoadDataFromFile(string(lsmType.LsmPath))
	//if err != nil {
	//	t.Fatalf("Error loading data from disk: %v", err)
	//}
	key := []byte(generateRandomKey())
	//key := []byte("UDVGKnSAsp")

	value := generateRandomByteSlices(4, 6)
	for _, v := range value {
		fmt.Println(string(v))
	}
	err := db.SAdd(key, value)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Get ok")
	fmt.Println(string(key))
	fmt.Println(string(lsmType.LsmPath))
	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}
