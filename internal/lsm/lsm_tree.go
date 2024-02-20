package lsm

import "sync"

// LSMTree 结构定义了 LSM 树的基本结构
type LSMTree struct {
	mu               sync.RWMutex      // 用于保护内存表的读写
	activeMemTable   map[string]string // 活跃的内存表，用于接收新写入的数据
	readOnlyMemTable map[string]string // 只读的内存表，用于读取操作
	diskLevels       [][]string        // 磁盘级别，存储已经持久化的数据
}

// NewLSMTree 函数用于创建一个新的 LSM 树实例
func NewLSMTree() *LSMTree {
	return &LSMTree{
		activeMemTable:   make(map[string]string),
		readOnlyMemTable: make(map[string]string),
		diskLevels:       make([][]string, 0),
	}
}
