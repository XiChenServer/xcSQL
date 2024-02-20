package lsm

import (
	"sync"
)

// LSMTree 结构定义了 LSM 树的基本结构
type LSMTree struct {
	mu               sync.RWMutex  // 用于保护内存表的读写
	activeMemTable   *SkipList     // 活跃的内存表，跳表作为索引
	readOnlyMemTable *SkipList     // 只读的内存表，跳表作为索引
	diskLevels       [][]*SkipList // 磁盘级别，存储已经持久化的数据，每个层级有多个跳表
	maxActiveSize    uint32        // 活跃内存表的最大大小
	maxDiskTableSize uint32        // 磁盘表的最大大小
}

// 在LSM树中插入数据
func (lsm *LSMTree) Insert(key []byte, value *DataInfo) {
	lsm.mu.Lock()

	// 检查活跃内存表的大小是否达到最大值，若达到则将活跃表转换为只读表，并写入磁盘
	if lsm.activeMemTable.Size >= lsm.maxActiveSize {
		lsm.convertActiveToReadOnly()
		lsm.writeReadOnlyToDisk()
	}
	lsm.mu.Unlock()
	// 插入数据到活跃内存表
	lsm.activeMemTable.Insert(key, value)
}

// 将活跃内存表转换为只读表
func (lsm *LSMTree) convertActiveToReadOnly() {
	lsm.readOnlyMemTable = lsm.activeMemTable
	lsm.activeMemTable = NewSkipList(3) // 重新初始化活跃内存表
}

// 将只读表的内容写入磁盘的第一级别
func (lsm *LSMTree) writeReadOnlyToDisk() {
	level := 0
	if len(lsm.diskLevels) <= level {
		lsm.diskLevels = append(lsm.diskLevels, []*SkipList{})
	}
	if len(lsm.diskLevels[level]) == 0 {
		lsm.diskLevels[level] = append(lsm.diskLevels[level], NewSkipList(3))
	}
	// 将只读表的内容写入磁盘的第一级别
	// 实现略
}
