package database

import "sync"

type XcDB struct {
	// 读写锁，用于并发读写控制
	mu sync.RWMutex
}
