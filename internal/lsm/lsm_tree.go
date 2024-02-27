package lsm

import (
	"bytes"
	"errors"
	"sync"
)

// LevelInfo 表示 LSM 树中的一个层级，包含该层级的跳表集合
type LevelInfo struct {
	SkipLists             []*SkipList // 该层级的跳表集合
	SkipListCount         uint16      // 该层级的跳表数量
	LevelMaxKey           []byte      // 该层级的最大键
	LevelMinKey           []byte      // 该层级的最小键
	LevelMaxSkipListCount uint16
}

// LSMTree 结构定义了 LSM 树的基本结构
type LSMTree struct {
	LsmPath              []byte       // 存储路径
	mu                   sync.RWMutex // 用于保护内存表的读写
	activeMemTable       *SkipList    // 活跃的内存表，跳表作为索引
	readOnlyMemTable     *SkipList    // 只读的内存表，跳表作为索引
	diskLevels           []*LevelInfo // 磁盘级别，存储已经持久化的数据，每个层级有多个跳表
	maxActiveSize        uint32       // 活跃内存表的最大大小
	maxDiskTableSize     uint32       // 磁盘表的最大大小
	maxSkipLists         uint16       // 每个层级的最大跳表数量
	maxDiskLevels        uint16
	writeToDiskWaitGroup sync.RWMutex  // 用于等待将只读表写入磁盘的协程完成
	writeToDiskChan      chan struct{} // 添加一个通道来控制磁盘写入的并发数

}

// 初始化 LSMTree
func NewLSMTree(maxActiveSize, maxDiskTableSize uint32) *LSMTree {
	maxSkipLists := uint16(10) // 第一个层级的跳表数量
	maxDiskLevels := uint16(7) // 最多的磁盘层级数量

	tree := &LSMTree{
		LsmPath:          []byte("../../data/testdata/lsm_tree/test1.txt"),
		activeMemTable:   NewSkipList(16),
		readOnlyMemTable: NewSkipList(16),
		diskLevels:       make([]*LevelInfo, maxDiskLevels),
		maxActiveSize:    maxActiveSize,
		maxDiskTableSize: maxDiskTableSize,
		maxSkipLists:     maxSkipLists,
		maxDiskLevels:    maxDiskLevels,
		writeToDiskChan:  make(chan struct{}, 1), // 初始化 writeToDiskChan，缓冲大小为 1，表示同时只能有一个磁盘写入操作
	}

	// 初始化每个层级的跳表数量
	skipLists := maxSkipLists
	for i := uint16(0); i < maxDiskLevels; i++ {
		// 为每个层级的 SkipLists 切片预分配空间
		//skipListSlice := make([]*SkipList, skipLists)
		tree.diskLevels[i] = &LevelInfo{
			SkipLists:             nil,
			SkipListCount:         0,
			LevelMaxKey:           []byte{}, // 使用空的字节数组表示最大键的缺失
			LevelMinKey:           []byte{}, // 使用空的字节数组表示最小键的缺失
			LevelMaxSkipListCount: skipLists,
		}
		skipLists *= 10 // 每个层级的跳表数量按4的幂级增加
	}

	return tree
}

// InsertAndMoveDown 方法用于插入数据到活跃内存表并执行跳表移动操作
func (lsm *LSMTree) Insert(key []byte, value *DataInfo) error {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	// 检查活跃内存表的大小是否达到最大值，若达到则将活跃表转换为只读表，并写入磁盘
	if lsm.activeMemTable.Size >= lsm.maxActiveSize {
		lsm.convertActiveToReadOnly()
		// 存储只读表到 LSM 树的最开始的层
		err := lsm.storeReadOnlyToFirstLevel(lsm.readOnlyMemTable)
		if err != nil {
			return err
		}
		lsm.readOnlyMemTable = NewSkipList(16) // 重新初始化只读内存表
	}

	// 插入数据到活跃内存表
	// 在插入时创建新的键值对副本，确保每个跳表保存的是独立的数据
	valueCopy := &DataInfo{
		DataMeta:        value.DataMeta,
		StorageLocation: value.StorageLocation,
	}
	lsm.activeMemTable.InsertInOrder(key, valueCopy)
	return nil
}

// Close 方法用于关闭 writeToDiskChan 通道
func (lsm *LSMTree) Close() {
	close(lsm.writeToDiskChan)
}

// InsertAndMoveDown 方法用于插入数据到活跃内存表并执行跳表移动操作
func (lsm *LSMTree) Get(key []byte) (*DataInfo, error) {
	lsm.mu.RLock()
	defer lsm.mu.RUnlock()
	// 在活跃内存表中查找与传入的键相同的项
	for node := lsm.activeMemTable.Head.Next[0]; node != nil; node = node.Next[0] {
		if bytes.Equal(node.Key, key) {
			return node.DataInfo, nil
		}
	}
	for node := lsm.readOnlyMemTable.Head.Next[0]; node != nil; node = node.Next[0] {
		if bytes.Equal(node.Key, key) {
			return node.DataInfo, nil
		}
	}
	for node := lsm.readOnlyMemTable.Head.Next[0]; node != nil; node = node.Next[0] {
		if bytes.Equal(node.Key, key) {
			return node.DataInfo, nil
		}
	}
	// 如果在活跃内存表中找不到，依次在其他层级中查找
	for _, level := range lsm.diskLevels {
		for _, skipList := range level.SkipLists {
			for node := skipList.Head.Next[0]; node != nil; node = node.Next[0] {
				if bytes.Equal(node.Key, key) {
					return node.DataInfo, nil
				}
			}
		}
	}
	err := errors.New("don`t find data")
	return nil, err // 如果找不到与传入的键相同的项，则返回 nil
}
