package lsm

import (
	"math/rand"
	"sync"
)

// Level 表示 LSM 树中的一个层级，包含该层级的跳表集合
type Level struct {
	SkipLists     []*SkipList // 该层级的跳表集合
	SkipListCount uint16      // 该层级的跳表数量
	LevelMaxKey   []byte      // 该层级的最大键
	LevelMinKey   []byte      // 该层级的最小键
}

// LSMTree 结构定义了 LSM 树的基本结构
type LSMTree struct {
	mu               sync.RWMutex // 用于保护内存表的读写
	activeMemTable   *SkipList    // 活跃的内存表，跳表作为索引
	readOnlyMemTable *SkipList    // 只读的内存表，跳表作为索引
	diskLevels       []*Level     // 磁盘级别，存储已经持久化的数据，每个层级有多个跳表
	maxActiveSize    uint32       // 活跃内存表的最大大小
	maxDiskTableSize uint32       // 磁盘表的最大大小
	maxSkipLists     uint16       // 每个层级的最大跳表数量
	maxDiskLevels    uint16
}

// 初始化 LSMTree
func NewLSMTree(maxActiveSize, maxDiskTableSize uint32) *LSMTree {
	maxSkipLists := uint16(4)  // 第一个层级的跳表数量
	maxDiskLevels := uint16(7) // 最多的磁盘层级数量

	tree := &LSMTree{
		activeMemTable:   NewSkipList(16),
		readOnlyMemTable: NewSkipList(16),
		diskLevels:       make([]*Level, maxDiskLevels),
		maxActiveSize:    maxActiveSize,
		maxDiskTableSize: maxDiskTableSize,
		maxSkipLists:     maxSkipLists,
		maxDiskLevels:    maxDiskLevels,
	}

	// 初始化每个层级的跳表数量
	skipLists := maxSkipLists
	for i := uint16(0); i < maxDiskLevels; i++ {
		tree.diskLevels[i] = &Level{
			SkipLists:     make([]*SkipList, skipLists),
			SkipListCount: skipLists,
		}
		skipLists *= 4 // 每个层级的跳表数量按4的幂级增加
	}

	return tree
}

// 在LSM树中插入数据
func (lsm *LSMTree) Insert(key []byte, value *DataInfo) {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	// 检查活跃内存表的大小是否达到最大值，若达到则将活跃表转换为只读表，并写入磁盘
	if lsm.activeMemTable.Size >= lsm.maxActiveSize {
		lsm.convertActiveToReadOnly()
		go lsm.writeReadOnlyToDisk() // 启动一个 goroutine 将只读表写入磁盘
	}

	// 插入数据到活跃内存表
	lsm.activeMemTable.Insert(key, value)
}

// 将活跃内存表转换为只读表
func (lsm *LSMTree) convertActiveToReadOnly() {
	lsm.readOnlyMemTable = lsm.activeMemTable
	lsm.activeMemTable = NewSkipList(16) // 重新初始化活跃内存表
}
func (lsm *LSMTree) writeReadOnlyToDisk() {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	// 存储只读表到第一层
	lsm.storeReadOnlyToLevel(lsm.readOnlyMemTable, 0)

	// 清空只读内存表
	lsm.readOnlyMemTable = nil
}

func (lsm *LSMTree) storeReadOnlyToLevel(skipList *SkipList, levelIndex int) {
	// 检查指定层是否已满
	if lsm.diskLevels[levelIndex].SkipListCount >= lsm.maxSkipLists {
		// 如果已满，先将该层的某个表移动到下一层
		lsm.moveSkipListDown(levelIndex)
	}

	// 将只读表存储到指定层
	lsm.diskLevels[levelIndex].SkipLists = append(lsm.diskLevels[levelIndex].SkipLists, skipList)
	lsm.diskLevels[levelIndex].SkipListCount++
}

func (lsm *LSMTree) moveSkipListDown(levelIndex int) {
	// 随机选择一个表移动到下一层
	randomIndex := rand.Intn(int(lsm.maxSkipLists))
	selectedSkipList := lsm.diskLevels[levelIndex].SkipLists[randomIndex]

	// 存储选定的跳表到下一层
	nextLevelIndex := levelIndex + 1
	lsm.storeReadOnlyToLevel(selectedSkipList, nextLevelIndex)

	// 从当前层中移除选定的跳表
	lsm.diskLevels[levelIndex].SkipLists = append(lsm.diskLevels[levelIndex].SkipLists[:randomIndex], lsm.diskLevels[levelIndex].SkipLists[randomIndex+1:]...)
}

// 检查当前层级是否能够容纳新的跳表
func (lsm *LSMTree) canCurrentLevelAccommodate() bool {
	level := len(lsm.diskLevels) // 获取当前磁盘级别的索引
	return lsm.diskLevels[level-1].SkipListCount < lsm.maxSkipLists
}

// 将只读表存储在指定的磁盘层级
func (lsm *LSMTree) storeReadOnlyMemTable(level int) {
	lsm.diskLevels[level-1].SkipLists = append(lsm.diskLevels[level-1].SkipLists, lsm.readOnlyMemTable)
	lsm.diskLevels[level-1].SkipListCount++
}

// 将某个跳表存储到下一层级
func (lsm *LSMTree) storeSelectedSkipListToNextLevel(level int) {
	// 从当前层级的跳表中随机选择一个
	randomIndex := rand.Intn(int(lsm.diskLevels[level-1].SkipListCount))
	selectedSkipList := lsm.diskLevels[level-1].SkipLists[randomIndex]

	// 将选定的跳表存储到下一层级
	nextLevel := level
	lsm.diskLevels[nextLevel].SkipLists = append(lsm.diskLevels[nextLevel].SkipLists, selectedSkipList)
	lsm.diskLevels[nextLevel].SkipListCount++

	// 从当前层级的跳表集合中移除选定的跳表
	lsm.diskLevels[level-1].SkipLists[randomIndex] = lsm.diskLevels[level-1].SkipLists[lsm.diskLevels[level-1].SkipListCount-1]
	lsm.diskLevels[level-1].SkipLists = lsm.diskLevels[level-1].SkipLists[:lsm.diskLevels[level-1].SkipListCount-1]

	// 更新当前层级的跳表数量
	lsm.diskLevels[level-1].SkipListCount--
}
