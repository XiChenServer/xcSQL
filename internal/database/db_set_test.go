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
	db := NewXcDB("")
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

// 简单的测试数据可以取出
func TestDB_SMembers(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB("")
	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
	lsmMap := *db.Lsm
	lsmType := lsmMap[model.XCDB_Set]
	//key := []byte(generateRandomKey())
	////key := []byte("UDVGKnSAsp")
	//
	//value := generateRandomByteSlices(4, 6)
	//for _, v := range value {
	//	fmt.Println(string(v))
	//}
	key := []byte("GHz2eStnJL")
	data, err := db.SMembers(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range data {
		fmt.Println(string(v))
	}
	fmt.Println(data)
	fmt.Println("Get ok")
	//fmt.Println(string(key))
	fmt.Println(string(lsmType.LsmPath))
	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}

// 简单的测试数据可以移除
func TestDB_SRem(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB("")
	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
	lsmMap := *db.Lsm
	lsmType := lsmMap[model.XCDB_Set]
	//key := []byte(generateRandomKey())
	////key := []byte("UDVGKnSAsp")
	//
	//value := generateRandomByteSlices(4, 6)
	//for _, v := range value {
	//	fmt.Println(string(v))
	//}
	key := []byte("GHz2eStnJL")
	value := [][]byte{[]byte("kpq"), []byte("B09")}
	err := db.SRem(key, value)
	if err != nil {
		fmt.Println(err)
		return
	}
	//for _, v := range data {
	//	fmt.Println(string(v))
	//}
	//fmt.Println(data)
	fmt.Println("Get ok")
	//fmt.Println(string(key))
	fmt.Println(string(lsmType.LsmPath))
	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}

// 简单的测试数据可以移除
func TestDB_SIsMember(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB("")
	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
	lsmMap := *db.Lsm
	lsmType := lsmMap[model.XCDB_Set]
	//key := []byte(generateRandomKey())
	////key := []byte("UDVGKnSAsp")
	//
	//value := generateRandomByteSlices(4, 6)
	//for _, v := range value {
	//	fmt.Println(string(v))
	//}
	key := []byte("GHz2eStnJL")
	value := []byte("kQR4h")
	flag, err := db.SIsMember(key, value)
	fmt.Println(flag)
	if err != nil {
		fmt.Println(err)
		return
	}
	//for _, v := range data {
	//	fmt.Println(string(v))
	//}
	//fmt.Println(data)
	fmt.Println("Get ok")
	//fmt.Println(string(key))
	fmt.Println(string(lsmType.LsmPath))
	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}

// 简单的测试数据可以移除
func TestDB_SCard(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB("")
	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
	lsmMap := *db.Lsm
	lsmType := lsmMap[model.XCDB_Set]
	//key := []byte(generateRandomKey())
	////key := []byte("UDVGKnSAsp")
	//
	//value := generateRandomByteSlices(4, 6)
	//for _, v := range value {
	//	fmt.Println(string(v))
	//}
	key := []byte("GHz2eStnJL")
	//value := []byte("kQR4h")
	flag, err := db.SCard(key)
	fmt.Println(flag)
	if err != nil {
		fmt.Println(err)
		return
	}
	//for _, v := range data {
	//	fmt.Println(string(v))
	//}
	//fmt.Println(data)
	fmt.Println("Get ok")
	//fmt.Println(string(key))
	fmt.Println(string(lsmType.LsmPath))
	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
}
