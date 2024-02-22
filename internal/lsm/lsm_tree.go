package lsm

import (
	"SQL/internal/database"
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
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
	mu               sync.RWMutex // 用于保护内存表的读写
	activeMemTable   *SkipList    // 活跃的内存表，跳表作为索引
	readOnlyMemTable *SkipList    // 只读的内存表，跳表作为索引
	diskLevels       []*LevelInfo // 磁盘级别，存储已经持久化的数据，每个层级有多个跳表
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
		diskLevels:       make([]*LevelInfo, maxDiskLevels),
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
		tree.diskLevels[i] = &LevelInfo{
			SkipLists:             skipListSlice,
			SkipListCount:         0,
			LevelMaxSkipListCount: skipLists,
		}
		skipLists *= 10 // 每个层级的跳表数量按4的幂级增加
	}

	return tree
}
func (lsm *LSMTree) Insert(key []byte, value *DataInfo) {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	// 检查活跃内存表的大小是否达到最大值，若达到则将活跃表转换为只读表，并写入磁盘
	if lsm.activeMemTable.Size >= lsm.maxActiveSize {
		lsm.convertActiveToReadOnly()
		go lsm.writeReadOnlyToDisk() // 启动一个 goroutine 将只读表写入磁盘
	}

	// 插入数据到活跃内存表
	// 在插入时创建新的键值对副本，确保每个跳表保存的是独立的数据
	valueCopy := &DataInfo{
		DataMeta:        value.DataMeta,
		StorageLocation: value.StorageLocation,
	}
	lsm.activeMemTable.InsertInOrder(key, valueCopy)
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
			lsm.diskLevels[levelIndex].SkipLists[lsm.diskLevels[levelIndex].SkipListCount] = newSkipList
			lsm.diskLevels[levelIndex].SkipListCount++

			// 更新层级的最大和最小键
			lsm.updateLevelMinMaxKeys(lsm.diskLevels[levelIndex], newSkipList)

			return
		}
	}
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

		// 更新当前层级的最大和最小键
		lsm.updateLevelMinMaxKeys(lsm.diskLevels[levelIndex], selectedSkipList)
		//lsm.updateLevelMinMaxKeys(lsm.diskLevels[nextLevelIndex], selectedSkipList)

		// 从当前层级的跳表中移除选定的跳表
		lsm.diskLevels[levelIndex].SkipLists[randomIndex] = lsm.diskLevels[levelIndex].SkipLists[skipListCount-1]
		lsm.diskLevels[levelIndex].SkipLists[skipListCount-1] = nil // 避免内存泄漏
		lsm.diskLevels[levelIndex].SkipListCount--

		return true
	}

	// 如果下一层已满，则递归调用 moveSkipListDown 函数，尝试将跳表移动到更下一层
	return lsm.moveSkipListDown(nextLevelIndex)
}
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

				// 打印跳表的最大和最小键
				writer.WriteString(fmt.Sprintf("SkipList %d, MaxKey: %s, MinKey: %s\n", skipListIndex, string(skipList.SkipListInfo.MaxKey), string(skipList.SkipListInfo.MinKey)))
			}
		}

		// 打印当前层级的最大和最小键
		writer.WriteString(fmt.Sprintf("Level %d, MaxKey: %s, MinKey: %s\n", levelIndex, string(level.LevelMaxKey), string(level.LevelMinKey)))
	}

	return nil
}

// 在程序退出时将活跃表保存到磁盘
func (lsm *LSMTree) SaveActiveToDiskOnExit() {
	lsm.readOnlyMemTable = lsm.activeMemTable
	// 在程序退出时保存活跃表到磁盘
	defer lsm.writeReadOnlyToDisk()
}
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

// LoadDataFromFile 从文件加载数据到 LSM 树中
func (lsm *LSMTree) LoadDataFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Level") {
			// 处理层级信息
			// 这里可以根据需要解析层级信息
			continue
		}
		if strings.HasPrefix(line, "SkipList") {
			// 处理跳表信息
			// 这里可以根据需要解析跳表信息
			continue
		}

		// 解析键值对信息
		keyValue := strings.Split(line, ", ")
		if len(keyValue) != 4 {
			return fmt.Errorf("invalid data format: %s", line)
		}

		key := []byte(strings.Split(keyValue[0], ": ")[1])
		value := []byte(strings.Split(keyValue[1], ": ")[1])
		extra := []byte(strings.Split(keyValue[2], ": ")[1])
		ttlStr := strings.Split(keyValue[3], ": ")[1]
		ttl, err := time.ParseDuration(ttlStr)
		if err != nil {
			return fmt.Errorf("failed to parse TTL: %v", err)
		}

		// 创建 DataInfo 对象
		data := DataInfo{
			DataMeta: database.DataMeta{
				Key:       key,
				Value:     value,
				Extra:     extra,
				KeySize:   uint32(len(key)),
				ValueSize: uint32(len(value)),
				ExtraSize: uint32(len(extra)),
				TTL:       ttl,
			},
		}

		// 将数据插入到 LSM 树中
		lsm.Insert(key, &data)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
