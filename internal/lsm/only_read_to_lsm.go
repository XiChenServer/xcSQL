package lsm

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
)

// // 保证只读的表存到lsm磁盘的第一层
//
//	func (lsm *LSMTree) storeReadOnlyToFirstLevel(skipList *SkipList) {
//		// 遍历磁盘级别，为每个层级创建新的跳表实例并复制数据
//		for levelIndex := 0; levelIndex < len(lsm.diskLevels); levelIndex++ {
//			// 创建新的只读表副本
//			readOnlyCopy := NewSkipList(16)
//			skipList.ForEach(func(key []byte, value *DataInfo) bool {
//				valueCopy := &DataInfo{
//					DataMeta:        value.DataMeta,
//					StorageLocation: value.StorageLocation,
//				}
//				readOnlyCopy.InsertInOrder(key, valueCopy)
//				return true
//			})
//
//			// 创建新的跳表实例
//			newSkipList := NewSkipList(16)
//
//			// 遍历只读表副本中的所有键值对，并插入到新的跳表中
//			readOnlyCopy.ForEach(func(key []byte, value *DataInfo) bool {
//				valueCopy := &DataInfo{
//					DataMeta:        value.DataMeta,
//					StorageLocation: value.StorageLocation,
//				}
//				newSkipList.InsertInOrder(key, valueCopy)
//				return true
//			})
//
//			// 检查当前层级是否已满
//			if lsm.diskLevels[levelIndex].SkipListCount < lsm.diskLevels[levelIndex].LevelMaxSkipListCount {
//				// 如果当前层级未满，则将新的跳表实例存储到该层级
//				//lsm.keepLsmLevelOrderly(levelIndex, newSkipList)
//				lsm.diskLevels[levelIndex].SkipLists[lsm.diskLevels[levelIndex].SkipListCount] = newSkipList
//				lsm.diskLevels[levelIndex].SkipListCount++
//
//				// 更新层级的最大和最小键
//				lsm.updateLevelMinMaxKeys(lsm.diskLevels[levelIndex], newSkipList)
//
//				return
//			}
//		}
//	}
//
// // 移动表到下一层，是一个递归的操作
//
//	func (lsm *LSMTree) moveSkipListDown(levelIndex int) {
//		// 如果当前层级的跳表数量为 0，则无法移动跳表到下一层
//		if lsm.diskLevels[levelIndex].SkipListCount == 0 {
//			return
//		}
//
//		// 随机选择一个表移动到下一层级
//		randomIndex := rand.Intn(int(lsm.diskLevels[levelIndex].SkipListCount))
//		selectedSkipList := lsm.diskLevels[levelIndex].SkipLists[randomIndex]
//
//		// 存储选定的跳表到下一层级
//		nextLevelIndex := levelIndex + 1
//
//		// 如果下一层已满，则递归调用移动操作，尝试将跳表移动到更下一层
//		if nextLevelIndex < len(lsm.diskLevels) && lsm.diskLevels[nextLevelIndex].SkipListCount >= lsm.diskLevels[nextLevelIndex].LevelMaxSkipListCount {
//			lsm.moveSkipListDown(nextLevelIndex)
//		}
//
//		// 检查下一层是否已满
//		if nextLevelIndex < len(lsm.diskLevels) && lsm.diskLevels[nextLevelIndex].SkipListCount < lsm.diskLevels[nextLevelIndex].LevelMaxSkipListCount {
//			// 将选定的跳表存储到下一层
//			lsm.diskLevels[nextLevelIndex].SkipLists[lsm.diskLevels[nextLevelIndex].SkipListCount] = selectedSkipList
//			//lsm.keepLsmLevelOrderly(levelIndex, selectedSkipList)
//			lsm.diskLevels[nextLevelIndex].SkipListCount++
//			// 更新层级的最大和最小键
//			lsm.updateLevelMinMaxKeys(lsm.diskLevels[levelIndex], selectedSkipList)
//
//			// 删除当前层级中选定的跳表
//			lsm.deleteSkipList(levelIndex, randomIndex)
//		}
//	}
//
// 将活跃内存表转换为只读表
func (lsm *LSMTree) convertActiveToReadOnly() {
	lsm.readOnlyMemTable = lsm.activeMemTable
	lsm.activeMemTable = NewSkipList(16) // 重新初始化活跃内存表
}

// 更新一个层中键的最大和最小的问题
func (lsm *LSMTree) updateLevelMinMaxKeys(currentLevel *LevelInfo, selectedSkipList *SkipList) {
	// 检查传入参数的有效性
	if currentLevel == nil || selectedSkipList == nil || selectedSkipList.SkipListInfo == nil {
		fmt.Println("Invalid currentLevel or selectedSkipList.")
		return
	}

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

func (lsm *LSMTree) storeReadOnlyToFirstLevel(skipList *SkipList) {
	// 如果第一层还有空间，则直接存储到第一层
	if lsm.diskLevels[0].SkipListCount >= lsm.diskLevels[0].LevelMaxSkipListCount {
		// 如果第一层已满，则随机选择一个跳表存储到下一层
		randomIndex := rand.Intn(int(lsm.diskLevels[0].SkipListCount))
		selectedSkipList := lsm.diskLevels[0].SkipLists[randomIndex]
		lsm.moveSkipListDown(0, randomIndex, selectedSkipList)

	}
	lsm.keepLsmLevelOrderly(0, skipList)

	//lsm.diskLevels[0].SkipLists = append(lsm.diskLevels[0].SkipLists, skipList)
	//lsm.diskLevels[0].SkipListCount++
	//lsm.updateLevelMinMaxKeys(lsm.diskLevels[0], skipList)
	return
}

// 将跳表从一个级别移动到下一个较低级别，并删除原来的位置
func (lsm *LSMTree) moveSkipListDown(levelIndex, skipListIndex int, skipList *SkipList) {
	// 如果下一层有空间，则存储到下一层
	if levelIndex+1 < len(lsm.diskLevels) && lsm.diskLevels[levelIndex+1].SkipListCount >= lsm.diskLevels[levelIndex+1].LevelMaxSkipListCount {
		// 如果下一层也满了，则随机选择一个跳表存储到下一层
		randomIndex := rand.Intn(int(lsm.diskLevels[levelIndex+1].SkipListCount))
		selectedSkipList := lsm.diskLevels[levelIndex+1].SkipLists[randomIndex]
		lsm.moveSkipListDown(levelIndex+1, randomIndex, selectedSkipList)
	}
	// 如果下一层有空间，则将选定的跳表存储到下一层

	lsm.keepLsmLevelOrderly(levelIndex+1, skipList)

	//lsm.diskLevels[levelIndex+1].SkipLists = append(lsm.diskLevels[levelIndex+1].SkipLists, skipList)
	//lsm.diskLevels[levelIndex+1].SkipListCount++
	//lsm.updateLevelMinMaxKeys(lsm.diskLevels[levelIndex+1], skipList)
	//删除第一层的原始位置
	lsm.deleteSkipList(levelIndex, skipListIndex)

}

// 删除指定层级的跳表
func (lsm *LSMTree) deleteSkipList(levelIndex, skipListIndex int) {
	// 检查待删除的索引是否有效
	if skipListIndex < 0 || skipListIndex >= len(lsm.diskLevels[levelIndex].SkipLists) {
		fmt.Println("Invalid skipListIndex.")
		return
	}

	// 将待删除的跳表从切片中移除
	skipLists := lsm.diskLevels[levelIndex].SkipLists
	copy(skipLists[skipListIndex:], skipLists[skipListIndex+1:])
	lsm.diskLevels[levelIndex].SkipLists = skipLists[:len(skipLists)-1]

	lsm.diskLevels[levelIndex].SkipListCount--
}

func (lsm *LSMTree) sortFirstLevelSkipLists() {
	sort.SliceStable(lsm.diskLevels[0].SkipLists, func(i, j int) bool {
		return bytes.Compare(lsm.diskLevels[0].SkipLists[i].SkipListInfo.MinKey, lsm.diskLevels[0].SkipLists[j].SkipListInfo.MinKey) < 0
	})
}
