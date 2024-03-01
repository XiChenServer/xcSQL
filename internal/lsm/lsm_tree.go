package lsm

import (
	"SQL/internal/model"
	"bytes"
	"errors"
	"fmt"
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
func NewLSMTree(maxActiveSize, maxDiskTableSize uint32, Type uint16) *LSMTree {
	maxSkipLists := uint16(10) // 第一个层级的跳表数量
	maxDiskLevels := uint16(7) // 最多的磁盘层级数量
	var typeName string
	if Type == model.String {
		typeName = "String"
	} else if Type == model.List {
		typeName = "List"
	}
	tree := &LSMTree{
		LsmPath:          []byte(("../../data/testdata/lsm_tree/") + typeName + ("/test1.gz")),
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
	//node := lsm.activeMemTable.Head
	//
	//for node != nil {
	//	fmt.Println("fd", string(node.Key))
	//	node = node.Next[0]
	//}
	// 插入数据到活跃内存表
	// 在插入时创建新的键值对副本，确保每个跳表保存的是独立的数据
	valueCopy := &DataInfo{
		DataMeta:        value.DataMeta,
		StorageLocation: value.StorageLocation,
	}
	lsm.activeMemTable.InsertInOrder(key, valueCopy)
	return nil
}

// InsertAndMoveDown 方法用于插入数据到活跃内存表并执行跳表移动操作
func (lsm *LSMTree) Insert1(key []byte, value *DataInfo) error {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()
	//
	// 检查活跃内存表的大小是否达到最大值，若达到则将活跃表转换为只读表，并写入磁盘
	if lsm.activeMemTable.Size >= lsm.maxActiveSize {
		lsm.convertActiveToReadOnly()
		//// 存储只读表到 LSM 树的最开始的层
		//err := lsm.storeReadOnlyToFirstLevel(lsm.readOnlyMemTable)
		//if err != nil {
		//	return err
		//}
		lsm.readOnlyMemTable = NewSkipList(16) // 重新初始化只读内存表

	}

	// 插入数据到活跃内存表
	// 在插入时创建新的键值对副本，确保每个跳表保存的是独立的数据
	valueCopy := &DataInfo{
		DataMeta:        value.DataMeta,
		StorageLocation: value.StorageLocation,
	}
	lsm.activeMemTable.InsertInOrder(key, valueCopy)

	node := lsm.activeMemTable.Head

	for node != nil {
		fmt.Println("fd", string(node.Key))
		node = node.Next[0]
	}
	return nil
}

// Close 方法用于关闭 writeToDiskChan 通道
func (lsm *LSMTree) Close() {
	close(lsm.writeToDiskChan)
}

// 写一个函数用于确定在哪个层级和跳表中进行查找
func (lsm *LSMTree) determineSearchRange(key []byte) (*LevelInfo, *SkipList) {
	// 在活跃内存表中查找
	if bytes.Compare(key, lsm.activeMemTable.getMinKey()) >= 0 && bytes.Compare(key, lsm.activeMemTable.getMaxKey()) <= 0 {
		fmt.Println("ff", string(key), string(lsm.activeMemTable.getMinKey()), string(lsm.activeMemTable.getMaxKey()))
		return nil, lsm.activeMemTable
	}

	// 在只读内存表中查找
	if bytes.Compare(key, lsm.readOnlyMemTable.getMinKey()) >= 0 && bytes.Compare(key, lsm.readOnlyMemTable.getMaxKey()) <= 0 {
		return nil, lsm.readOnlyMemTable
	}

	// 在磁盘层级中查找
	for _, level := range lsm.diskLevels {
		// 首先检查目标键是否在当前层级的键范围内
		if bytes.Compare(key, level.LevelMinKey) >= 0 && bytes.Compare(key, level.LevelMaxKey) <= 0 {
			// 使用二分查找确定目标键所在的跳表
			low, high := 0, len(level.SkipLists)-1
			for low <= high {
				mid := (low + high) / 2
				midKey := level.SkipLists[mid].getMaxKey()
				if bytes.Compare(midKey, key) < 0 {
					low = mid + 1
				} else if bytes.Compare(midKey, key) > 0 {
					high = mid - 1
				} else {
					// 找到目标键所在的跳表
					return level, level.SkipLists[mid]
				}
			}
		}
	}

	// 如果未找到目标键所在的跳表，则返回 nil
	return nil, nil

}

// 修改 Get 函数，使用二分查找算法在确定的层级和跳表中进行查找
func (lsm *LSMTree) Get1(key []byte) (*DataInfo, error) {
	lsm.mu.RLock()
	defer lsm.mu.RUnlock()

	level, skipList := lsm.determineSearchRange(key)
	if level == nil && skipList == nil {
		return nil, errors.New("don't find data")
	}

	for node := skipList.Head.Next[0]; node != nil; node = node.Next[0] {
		fmt.Println(string(node.Key), string(node.DataInfo.Value), string(node.DataInfo.Key))
		if bytes.Equal(node.Key, key) {
			return node.DataInfo, nil
		}
	}
	//if curNode.Next[0] != nil && bytes.Equal(curNode.Next[0].Key, key) {
	//	return curNode.Next[0].DataInfo, nil
	//}

	return nil, errors.New("don't find data")
}

// 修改 Get 函数，先在活跃表、只读表、LSM 树的层级中的跳表中进行查找
func (lsm *LSMTree) Get(key []byte) (*DataInfo, error) {
	lsm.mu.RLock()
	defer lsm.mu.RUnlock()

	// 先在活跃表中查找
	if bytes.Compare(key, lsm.activeMemTable.getMinKey()) >= 0 && bytes.Compare(key, lsm.activeMemTable.getMaxKey()) <= 0 {
		return lsm.searchInSkipList(lsm.activeMemTable, key)
	}

	// 然后在只读表中查找
	if bytes.Compare(key, lsm.readOnlyMemTable.getMinKey()) >= 0 && bytes.Compare(key, lsm.readOnlyMemTable.getMaxKey()) <= 0 {
		return lsm.searchInSkipList(lsm.readOnlyMemTable, key)
	}

	// 最后在 LSM 树的层级中的跳表中进行查找
	for _, level := range lsm.diskLevels {
		if bytes.Compare(key, level.LevelMinKey) >= 0 && bytes.Compare(key, level.LevelMaxKey) <= 0 {
			for _, skipList := range level.SkipLists {
				if bytes.Compare(key, skipList.getMinKey()) >= 0 && bytes.Compare(key, skipList.getMaxKey()) <= 0 {
					return lsm.searchInSkipList(skipList, key)
				}
			}
		}
	}

	return nil, errors.New("data not found")
}

// 辅助函数：在跳表中查找数据
func (lsm *LSMTree) searchInSkipList(skipList *SkipList, key []byte) (*DataInfo, error) {
	for node := skipList.Head.Next[0]; node != nil; node = node.Next[0] {
		if bytes.Equal(node.Key, key) {
			return node.DataInfo, nil
		} else if bytes.Compare(node.Key, key) > 0 {
			break // 如果当前节点的键大于目标键，则跳出循环
		}
	}

	return nil, errors.New("data not found")
}

func (lsm *LSMTree) Printf() {
	node := lsm.activeMemTable.Head
	for node != nil {
		fmt.Println("1", string(node.Key))
		node = node.Next[0]
	}
}
