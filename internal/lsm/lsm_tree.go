package lsm

import (
	"bytes"
	"fmt"
	"math/rand"
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
		skipListSlice := make([]*SkipList, skipLists)
		tree.diskLevels[i] = &LevelInfo{
			SkipLists:             skipListSlice,
			SkipListCount:         0,
			LevelMaxKey:           []byte{}, // 使用空的字节数组表示最大键的缺失
			LevelMinKey:           []byte{}, // 使用空的字节数组表示最小键的缺失
			LevelMaxSkipListCount: skipLists,
		}
		skipLists *= 10 // 每个层级的跳表数量按4的幂级增加
	}

	return tree
}

//func (lsm *LSMTree) Insert(key []byte, value *DataInfo) {
//	lsm.mu.Lock()
//	defer lsm.mu.Unlock()
//
//	// 检查活跃内存表的大小是否达到最大值，若达到则将活跃表转换为只读表，并写入磁盘
//	if lsm.activeMemTable.Size >= lsm.maxActiveSize {
//		lsm.convertActiveToReadOnly()
//		// 使用锁来保证只有出现竞态的问题
//		lsm.writeToDiskWaitGroup.Lock()
//		lsm.writeReadOnlyToDisk()
//		lsm.writeToDiskWaitGroup.Unlock()
//		//lsm.writeToDiskChan <- struct{}{} // 发送信号到 writeToDiskChan 通道
//		lsm.activeMemTable = NewSkipList(16)
//	}
//	// 插入数据到活跃内存表
//	// 在插入时创建新的键值对副本，确保每个跳表保存的是独立的数据
//	valueCopy := &DataInfo{
//		DataMeta:        value.DataMeta,
//		StorageLocation: value.StorageLocation,
//	}
//	lsm.activeMemTable.InsertInOrder(key, valueCopy)
//}

// 将活跃内存表转换为只读表
func (lsm *LSMTree) convertActiveToReadOnly() {
	lsm.readOnlyMemTable = lsm.activeMemTable
	lsm.activeMemTable = NewSkipList(16) // 重新初始化活跃内存表
}

// 将只读表存到lsm的磁盘之中
func (lsm *LSMTree) writeReadOnlyToDisk() {

	// 存储只读表到第一层
	lsm.storeReadOnlyToFirstLevel(lsm.readOnlyMemTable)

	//// 清空只读内存表
	//lsm.readOnlyMemTable = nil
}

// InsertAndMoveDown 方法用于插入数据到活跃内存表并执行跳表移动操作
func (lsm *LSMTree) Insert(key []byte, value *DataInfo) {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	// 检查活跃内存表的大小是否达到最大值，若达到则将活跃表转换为只读表，并写入磁盘
	if lsm.activeMemTable.Size >= lsm.maxActiveSize {
		lsm.convertActiveToReadOnly()
		// 检查最开始的层是否已满，如果已满，则触发跳表移动操作
		if lsm.diskLevels[0].SkipListCount >= lsm.diskLevels[0].LevelMaxSkipListCount {
			lsm.moveSkipListDown(0)
		}
		// 存储只读表到 LSM 树的最开始的层
		lsm.storeReadOnlyToFirstLevel(lsm.readOnlyMemTable)
		lsm.readOnlyMemTable = NewSkipList(16) // 重新初始化只读内存表
	}

	// 插入数据到活跃内存表
	// 在插入时创建新的键值对副本，确保每个跳表保存的是独立的数据
	valueCopy := &DataInfo{
		DataMeta:        value.DataMeta,
		StorageLocation: value.StorageLocation,
	}
	lsm.activeMemTable.InsertInOrder(key, valueCopy)
}

// 保证只读的表存到lsm磁盘的第一层
func (lsm *LSMTree) storeReadOnlyToFirstLevel(skipList *SkipList) {
	// 遍历磁盘级别，为每个层级创建新的跳表实例并复制数据
	for levelIndex := 0; levelIndex < len(lsm.diskLevels); levelIndex++ {
		// 创建新的只读表副本
		readOnlyCopy := NewSkipList(16)
		skipList.ForEach(func(key []byte, value *DataInfo) bool {
			valueCopy := &DataInfo{
				DataMeta:        value.DataMeta,
				StorageLocation: value.StorageLocation,
			}
			readOnlyCopy.InsertInOrder(key, valueCopy)
			return true
		})

		// 创建新的跳表实例
		newSkipList := NewSkipList(16)

		// 遍历只读表副本中的所有键值对，并插入到新的跳表中
		readOnlyCopy.ForEach(func(key []byte, value *DataInfo) bool {
			valueCopy := &DataInfo{
				DataMeta:        value.DataMeta,
				StorageLocation: value.StorageLocation,
			}
			newSkipList.InsertInOrder(key, valueCopy)
			return true
		})

		// 检查当前层级是否已满
		if lsm.diskLevels[levelIndex].SkipListCount < lsm.diskLevels[levelIndex].LevelMaxSkipListCount {
			// 如果当前层级未满，则将新的跳表实例存储到该层级
			lsm.keepLsmLevelOrderly(levelIndex, newSkipList)
			//lsm.diskLevels[levelIndex].SkipLists[lsm.diskLevels[levelIndex].SkipListCount] = newSkipList
			lsm.diskLevels[levelIndex].SkipListCount++

			// 更新层级的最大和最小键
			lsm.updateLevelMinMaxKeys(lsm.diskLevels[levelIndex], newSkipList)

			return
		}
	}
}

// 在整个跳表进行插入的时候，保证lsm整个层的有序性
func (lsm *LSMTree) keepLsmLevelOrderly(levelIndex int, skipList *SkipList) {
	if lsm.diskLevels[levelIndex].SkipListCount == 0 {
		fmt.Println("1232")
		lsm.diskLevels[levelIndex].SkipLists = append(lsm.diskLevels[levelIndex].SkipLists, skipList)
		return
	}

	fmt.Println("Before calling keepLsmLevelOrderly:")
	fmt.Println("LevelIndex:", levelIndex)
	fmt.Println("LSMTree:", lsm)
	fmt.Println("DiskLevels:", lsm.diskLevels)

	// 添加对 skipList 是否为空的有效性检查
	if skipList == nil || skipList.SkipListInfo == nil {
		fmt.Println("Invalid skipList or skipListInfo is nil.")
		return
	}

	levelMinKey := lsm.diskLevels[levelIndex].LevelMinKey
	levelMaxKey := lsm.diskLevels[levelIndex].LevelMaxKey
	if bytes.Compare(levelMinKey, skipList.SkipListInfo.MaxKey) >= 0 {
		// 如果新跳表的最大键大于等于当前层级的最小键，直接将新跳表插入到当前层级的首位
		lsm.diskLevels[levelIndex].SkipLists = append([]*SkipList{skipList}, lsm.diskLevels[levelIndex].SkipLists...)
	} else if bytes.Compare(levelMaxKey, skipList.SkipListInfo.MinKey) <= 0 {
		// 如果新跳表的最小键小于等于当前层级的最大键，直接将新跳表插入到当前层级的末尾
		lsm.diskLevels[levelIndex].SkipLists = append(lsm.diskLevels[levelIndex].SkipLists, skipList)
	} else {
		// 否则，需要找到新跳表应该插入的位置，确保整个层级的有序性
		for i, existingSkipList := range lsm.diskLevels[levelIndex].SkipLists {
			if bytes.Compare(existingSkipList.SkipListInfo.MaxKey, skipList.SkipListInfo.MaxKey) > 0 {
				// 如果当前跳表的最大键大于新跳表的最大键，说明新跳表应该插入到当前位置的前面
				lsm.diskLevels[levelIndex].SkipLists = append(lsm.diskLevels[levelIndex].SkipLists[:i], append([]*SkipList{skipList}, lsm.diskLevels[levelIndex].SkipLists[i:]...)...)
				break
			}
		}
	}
}

// 移动表到下一层，是一个递归的操作
func (lsm *LSMTree) moveSkipListDown(levelIndex int) {
	// 如果当前层级的跳表数量为 0，则无法移动跳表到下一层
	if lsm.diskLevels[levelIndex].SkipListCount == 0 {
		return
	}

	// 随机选择一个表移动到下一层级
	randomIndex := rand.Intn(int(lsm.diskLevels[levelIndex].SkipListCount))
	selectedSkipList := lsm.diskLevels[levelIndex].SkipLists[randomIndex]

	// 存储选定的跳表到下一层级
	nextLevelIndex := levelIndex + 1

	// 如果下一层已满，则递归调用移动操作，尝试将跳表移动到更下一层
	if nextLevelIndex < len(lsm.diskLevels) && lsm.diskLevels[nextLevelIndex].SkipListCount >= lsm.diskLevels[nextLevelIndex].LevelMaxSkipListCount {
		lsm.moveSkipListDown(nextLevelIndex)
	}

	// 检查下一层是否已满
	if nextLevelIndex < len(lsm.diskLevels) && lsm.diskLevels[nextLevelIndex].SkipListCount < lsm.diskLevels[nextLevelIndex].LevelMaxSkipListCount {
		// 将选定的跳表存储到下一层
		lsm.diskLevels[nextLevelIndex].SkipLists[lsm.diskLevels[nextLevelIndex].SkipListCount] = selectedSkipList
		lsm.keepLsmLevelOrderly(levelIndex, selectedSkipList)
		//lsm.diskLevels[nextLevelIndex].SkipListCount++
		// 更新层级的最大和最小键
		lsm.updateLevelMinMaxKeys(lsm.diskLevels[levelIndex], selectedSkipList)

		// 删除当前层级中选定的跳表
		lsm.deleteSkipList(levelIndex, randomIndex)
	}
}

// 删除指定层级的跳表
func (lsm *LSMTree) deleteSkipList(levelIndex, skipListIndex int) {
	// 将要删除的跳表替换为最后一个跳表，并将计数减一
	lastIndex := int(lsm.diskLevels[levelIndex].SkipListCount) - 1
	lsm.diskLevels[levelIndex].SkipLists[skipListIndex] = lsm.diskLevels[levelIndex].SkipLists[lastIndex]
	lsm.diskLevels[levelIndex].SkipLists[lastIndex] = nil
	lsm.diskLevels[levelIndex].SkipListCount--
}

// Close 方法用于关闭 writeToDiskChan 通道
func (lsm *LSMTree) Close() {
	close(lsm.writeToDiskChan)
}

// 更新一个层中键的最大和最小的问题
func (lsm *LSMTree) updateLevelMinMaxKeys(currentLevel *LevelInfo, selectedSkipList *SkipList) {
	// 获取跳表的最小键和最大键
	minKey := selectedSkipList.SkipListInfo.MinKey
	maxKey := selectedSkipList.SkipListInfo.MaxKey

	// 如果跳表为空，则直接返回
	if minKey == nil || maxKey == nil {
		return
	}

	// 如果当前层级的最小键为空或者跳表的最小键小于当前层级的最小键，则更新最小键
	if len(currentLevel.LevelMinKey) == 0 || bytes.Compare(minKey, currentLevel.LevelMinKey) < 0 {
		currentLevel.LevelMinKey = minKey
	}

	// 如果当前层级的最大键为空或者跳表的最大键大于当前层级的最大键，则更新最大键
	if len(currentLevel.LevelMaxKey) == 0 || bytes.Compare(maxKey, currentLevel.LevelMaxKey) > 0 {
		currentLevel.LevelMaxKey = maxKey
	}
}
