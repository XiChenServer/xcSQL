package database

import (
	"SQL/internal/lsm"
	"SQL/internal/model"
	"SQL/internal/storage"
	"SQL/logs"
	"sync"
)

type XcDB struct {
	StorageManager *storage.StorageManager
	Lsm            *map[uint16]*lsm.LSMTree
	// 读写锁，用于并发读写控制
	Mu sync.RWMutex
}

func NewXcDB() *XcDB {
	var lsmMap = make(map[uint16]*lsm.LSMTree)
	// 启动一个协程来初始化字符串类型的LSM树
	//go func() {
	//	lsmString := lsm.NewLSMTree(16, 10000, model.String)
	//	// 在这里可以对 lsmString 进行操作，例如插入初始数据等
	//	lsmMap[model.String] = *lsmString
	//}()
	//
	//// 启动一个协程来初始化列表类型的LSM树
	//go func() {
	//	lsmList := lsm.NewLSMTree(16, 10000, model.List)
	//	// 在这里可以对 lsmList 进行操作，例如插入初始数据等
	//	lsmMap[model.List] = *lsmList
	//}()
	lsmString := lsm.NewLSMTree(16, 10000, model.String)
	lsmList := lsm.NewLSMTree(16, 10000, model.List)
	lsmMap[model.List] = lsmList
	lsmMap[model.String] = lsmString
	storageManager, err := storage.LoadStorageManager("../../data/testdata/lsm_tree/config.txt")

	if err != nil {
		storageManager, err = storage.NewStorageManager("../../data/testdata/string_test", 10*1024) // 1MB 文件大小限制
		logs.SugarLogger.Error("failed to create storage manager: %v", err)
	}
	return &XcDB{
		Lsm:            &lsmMap,
		StorageManager: storageManager,
		Mu:             sync.RWMutex{},
	}
}
