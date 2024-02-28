package database

import (
	"SQL/internal/lsm"
	"SQL/internal/storage"
	"SQL/logs"
	"sync"
)

type XcDB struct {
	StorageManager *storage.StorageManager
	Lsm            *lsm.LSMTree
	// 读写锁，用于并发读写控制
	Mu sync.RWMutex
}

func NewXcDB() *XcDB {

	lsm := lsm.NewLSMTree(16, 10000)
	storageManager, err := storage.LoadStorageManager("../../data/testdata/lsm_tree/config.txt")
	if err != nil {
		storageManager, err = storage.NewStorageManager("../../data/testdata/string_test", 10*1024) // 1MB 文件大小限制
		logs.SugarLogger.Error("failed to create storage manager: %v", err)
	}
	return &XcDB{
		Lsm:            lsm,
		StorageManager: storageManager,
		Mu:             sync.RWMutex{},
	}
}
