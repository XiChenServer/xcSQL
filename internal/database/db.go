package database

import (
	"SQL/internal/storage"
	"SQL/logs"
	"sync"
)

type XcDB struct {
	storageManager *storage.StorageManager
	// 读写锁，用于并发读写控制
	mu sync.RWMutex
}

func NewXcDB() *XcDB {
	storageManager, err := storage.NewStorageManager("../../data/testdata/string_test", 4*1024) // 1MB 文件大小限制
	if err != nil {
		logs.SugarLogger.Panicln("failed to create storage manager: %v", err)
	}
	return &XcDB{
		storageManager: storageManager,
		mu:             sync.RWMutex{},
	}
}
