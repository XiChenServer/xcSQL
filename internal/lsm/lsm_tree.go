package lsm

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
)

// Level 表示 LSM 树中的一个层级，包含该层级的跳表集合
type Level struct {
	SkipLists             []*SkipList // 该层级的跳表集合
	SkipListCount         uint16      // 该层级的跳表数量
	LevelMaxKey           []byte      // 该层级的最大键
	LevelMinKey           []byte      // 该层级的最小键
	LevelMaxSkipListCount uint16
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
	maxSkipLists := uint16(10) // 第一个层级的跳表数量
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
		// 为每个层级的 SkipLists 切片预分配空间
		skipListSlice := make([]*SkipList, skipLists)
		tree.diskLevels[i] = &Level{
			SkipLists:             skipListSlice,
			SkipListCount:         0,
			LevelMaxSkipListCount: skipLists,
		}
		skipLists *= 10 // 每个层级的跳表数量按4的幂级增加
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
	lsm.storeReadOnlyToFirstLevel(lsm.readOnlyMemTable)

	//// 清空只读内存表
	//lsm.readOnlyMemTable = nil
}
func (lsm *LSMTree) storeReadOnlyToFirstLevel(skipList *SkipList) {
	// 检查第一层是否已满
	if lsm.diskLevels[0].SkipListCount >= lsm.diskLevels[0].LevelMaxSkipListCount {
		// 如果第一层已满，先将第一层的某个表移动到第二层
		moved := lsm.moveSkipListDown(0)

		// 如果成功移动了跳表，则在第一层存储只读表
		if moved {
			lsm.diskLevels[0].SkipLists[lsm.diskLevels[0].SkipListCount] = skipList
			lsm.diskLevels[0].SkipListCount++
			return
		}
	}

	// 如果第一层未满，则直接在第一层存储只读表
	lsm.diskLevels[0].SkipLists[lsm.diskLevels[0].SkipListCount] = skipList
	lsm.diskLevels[0].SkipListCount++
}

func (lsm *LSMTree) moveSkipListDown(levelIndex int) bool {
	// 获取当前层级的跳表数量
	skipListCount := int(lsm.diskLevels[levelIndex].SkipListCount)

	// 如果当前层级的跳表数量为 0，则无法移动跳表到下一层
	if skipListCount == 0 {
		return false
	}

	// 随机选择一个表移动到下一层
	randomIndex := rand.Intn(skipListCount)
	selectedSkipList := lsm.diskLevels[levelIndex].SkipLists[randomIndex]

	// 存储选定的跳表到下一层级
	nextLevelIndex := levelIndex + 1

	// 检查下一层是否已满
	if nextLevelIndex < len(lsm.diskLevels) && lsm.diskLevels[nextLevelIndex].SkipListCount < lsm.diskLevels[nextLevelIndex].LevelMaxSkipListCount {
		// 将选定的跳表存储到下一层
		lsm.diskLevels[nextLevelIndex].SkipLists[lsm.diskLevels[nextLevelIndex].SkipListCount] = selectedSkipList
		lsm.diskLevels[nextLevelIndex].SkipListCount++

		// 从当前层级的跳表中移除选定的跳表
		lsm.diskLevels[levelIndex].SkipLists[randomIndex] = lsm.diskLevels[levelIndex].SkipLists[skipListCount-1]
		lsm.diskLevels[levelIndex].SkipLists[skipListCount-1] = nil // 避免内存泄漏
		lsm.diskLevels[levelIndex].SkipListCount--

		return true
	}

	// 如果下一层已满，则递归调用 moveSkipListDown 函数，尝试将跳表移动到更下一层
	return lsm.moveSkipListDown(nextLevelIndex)
}

// 将磁盘上的所有数据打印并保存到文件
func (lsm *LSMTree) PrintDiskDataToFile(filePath string) error {
	lsm.mu.RLock()
	defer lsm.mu.RUnlock()

	// 打开文件准备写入
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 bufio.Writer 提高写入性能
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 遍历每个层级的跳表进行查找
	for levelIndex, level := range lsm.diskLevels {
		for skipListIndex, skipList := range level.SkipLists {
			if skipList != nil {
				writer.WriteString(fmt.Sprintf("Level %d, SkipList %d:\n", levelIndex, skipListIndex))
				// 遍历跳表中的所有键值对并写入文件
				skipList.ForEach(func(key []byte, value *DataInfo) bool {
					line := fmt.Sprintf("Key: %s, Value: %s, Extra: %s, TTL: %s\n", string(key), string(value.Value), string(value.Extra), value.TTL.String())
					writer.WriteString(line)
					return true
				})
			}
		}
	}

	return nil

}

// 在程序退出时将活跃表保存到磁盘
func (lsm *LSMTree) SaveActiveToDiskOnExit() {
	lsm.readOnlyMemTable = lsm.activeMemTable
	// 在程序退出时保存活跃表到磁盘
	defer lsm.writeReadOnlyToDisk()
}
