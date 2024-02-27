package database

import (
	"SQL/internal/lsm"
	"SQL/internal/storage"
	"SQL/logs"
	"sync"
)

type XcDB struct {
	storageManager *storage.StorageManager
	lsm            *lsm.LSMTree
	// 读写锁，用于并发读写控制
	mu sync.RWMutex
}

func NewXcDB() *XcDB {

	lsm := lsm.NewLSMTree(16, 10000)
	storageManager, err := storage.LoadStorageManager("../../data/testdata/lsm_tree/config.txt")
	if err != nil {
		storageManager, err = storage.NewStorageManager("../../data/testdata/string_test", 4*1024) // 1MB 文件大小限制
		logs.SugarLogger.Error("failed to create storage manager: %v", err)
	}
	return &XcDB{
		lsm:            lsm,
		storageManager: storageManager,
		mu:             sync.RWMutex{},
	}
}
