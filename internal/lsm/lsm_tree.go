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

// NewLSMTree 创建一个新的 LSM 树
func NewLSMTree(maxActiveSize, maxDiskTableSize, initialLevelCapacity uint32) *LSMTree {
	lsm := &LSMTree{
		maxActiveSize:    maxActiveSize,
		maxDiskTableSize: maxDiskTableSize,
	}
	// 初始化活跃内存表
	lsm.activeMemTable = NewSkipList(16)
	// 初始化只读内存表
	lsm.readOnlyMemTable = NewSkipList(16)
	// 初始化磁盘级别
	lsm.diskLevels = []*Level{{SkipLists: make([]*SkipList, 0, initialLevelCapacity), SkipListCount: 0}}
	return lsm
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

// 将只读表的内容写入磁盘的最上一层
func (lsm *LSMTree) writeReadOnlyToDisk() {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	level := 0 // 获取当前磁盘级别的索引

	// 检查当前层级是否能够容纳新的跳表
	if lsm.diskLevels[level].SkipListCount >= lsm.maxSkipLists {
		// 如果最上面一层已满，找到第一个非满的层级
		for i := 0; i < len(lsm.diskLevels); i++ {
			if lsm.diskLevels[i].SkipListCount < lsm.maxSkipLists {
				level = i
				break
			}
		}

		// 如果所有层级都已满，则新增一层
		if lsm.diskLevels[level].SkipListCount >= lsm.maxSkipLists {
			lsm.addNewDiskLevel(level)
		}
	}

	// 存储只读表到当前或新增的层级
	lsm.diskLevels[level].SkipLists = append(lsm.diskLevels[level].SkipLists, lsm.readOnlyMemTable)
	lsm.diskLevels[level].SkipListCount++
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

// 增加一个新的磁盘层级
func (lsm *LSMTree) addNewDiskLevel(level int) {
	lsm.diskLevels = append(lsm.diskLevels, &Level{})
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
