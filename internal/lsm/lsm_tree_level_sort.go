package lsm

import (
	"bytes"
	"errors"
	"fmt"
)

// 在整个跳表进行插入的时候，保证LSM整个层的有序性
func (lsm *LSMTree) keepLsmLevelOrderly(levelIndex int, skipList *SkipList) error {
	// 检查传入参数的有效性
	if levelIndex < 0 || levelIndex >= int(lsm.maxDiskLevels) {
		fmt.Println("Invalid level index.")
		err := errors.New("Invalid skipListIndex.")
		return err
	}

	if skipList == nil || skipList.SkipListInfo == nil || skipList.Head == nil {
		fmt.Println("Invalid skipList or skipListInfo is nil.")
		err := errors.New("Invalid skipList or skipListInfo is nil.")
		return err
	}

	// 遍历已存在的跳表，找到新跳表应该插入的位置
	for i, existingSkipList := range lsm.diskLevels[levelIndex].SkipLists {
		if existingSkipList == nil || existingSkipList.SkipListInfo == nil {
			continue
		}

		// 比较新跳表的最大键和当前跳表的最大键
		if bytes.Compare(existingSkipList.SkipListInfo.MaxKey, skipList.SkipListInfo.MaxKey) > 0 {
			// 在当前位置插入新跳表
			lsm.diskLevels[levelIndex].SkipLists = append(lsm.diskLevels[levelIndex].SkipLists[:i], append([]*SkipList{skipList}, lsm.diskLevels[levelIndex].SkipLists[i:]...)...)
			lsm.diskLevels[levelIndex].SkipListCount++
			lsm.updateLevelMinMaxKeys(lsm.diskLevels[levelIndex], skipList)
			return nil
		}
	}

	// 如果新跳表的最大键大于等于当前所有跳表的最大键，则将新跳表追加到末尾
	lsm.diskLevels[levelIndex].SkipLists = append(lsm.diskLevels[levelIndex].SkipLists, skipList)
	lsm.diskLevels[levelIndex].SkipListCount++
	err := lsm.updateLevelMinMaxKeys(lsm.diskLevels[levelIndex], skipList)
	if err != nil {
		return err
	}
	return nil
}

// 检查lsm的哪一层中是否存在重叠跳表
func (lsm *LSMTree) hasOverlappingSkipLists(levelIndex int, skipList *SkipList) bool {
	// 获取只读内存表的最大键和最小键
	maxKey := skipList.getMaxKey()
	minKey := skipList.getMinKey()

	// 遍历所有层级的跳表，检查是否存在重叠
	for _, sl := range lsm.diskLevels[levelIndex].SkipLists {
		if sl == nil {
			continue
		}
		slMaxKey := sl.getMaxKey()
		slMinKey := sl.getMinKey()

		// 如果存在重叠，则返回 true
		if bytes.Compare(minKey, slMaxKey) < 0 && bytes.Compare(maxKey, slMinKey) > 0 {
			return true
		}
	}

	return false
}

// 合并只读内存表中的重叠跳表
func (lsm *LSMTree) mergeOverlappingSkipLists(levelIndex int, skipList *SkipList) {
	// 获取当前层级的跳表列表
	level := lsm.diskLevels[levelIndex]

	// 获取待插入跳表的最大键和最小键
	maxKey := skipList.getMaxKey()
	minKey := skipList.getMinKey()

	// 初始化一个切片，用于存储合并后的跳表列表
	mergedSkipLists := make([]*SkipList, 0)

	// 遍历当前层级的跳表，查找需要合并的跳表
	for _, sl := range level.SkipLists {
		if sl == nil {
			continue
		}
		slMaxKey := sl.getMaxKey()
		slMinKey := sl.getMinKey()

		// 如果存在重叠，则进行合并操作
		if bytes.Compare(minKey, slMaxKey) < 0 && bytes.Compare(maxKey, slMinKey) > 0 {
			// 合并跳表操作
			skipList = mergeSortedSkipLists(skipList, sl)
		} else {
			// 如果不存在重叠，则将原跳表保留在合并后的列表中
			mergedSkipLists = append(mergedSkipLists, sl)
		}
	}

	// 将合并后的跳表添加到列表中
	mergedSkipLists = append(mergedSkipLists, skipList)

	// 更新当前层级的跳表列表
	level.SkipLists = mergedSkipLists
}

// 合并两个有序跳表
func mergeSortedSkipLists(sl1, sl2 *SkipList) *SkipList {
	maxLevel := sl1.MaxLevel                // 假设两个跳表的层数相同，选择其中一个跳表的层数作为合并后跳表的层数
	mergedSkipList := NewSkipList(maxLevel) // 创建一个新的跳表作为合并后的结果
	node1 := sl1.Head.Next[0]               // 第一个跳表的头节点
	node2 := sl2.Head.Next[0]               // 第二个跳表的头节点

	// 双指针遍历两个跳表的节点
	for node1 != nil && node2 != nil {
		// 如果节点1的键小于节点2的键，则将节点1插入到合并后的跳表中
		if bytes.Compare(node1.Key, node2.Key) < 0 {
			mergedSkipList.InsertInOrder(node1.Key, node1.DataInfo)
			node1 = node1.Next[0] // 移动节点1的指针
		} else {
			// 否则，将节点2插入到合并后的跳表中
			mergedSkipList.InsertInOrder(node2.Key, node2.DataInfo)
			node2 = node2.Next[0] // 移动节点2的指针
		}
	}

	// 将剩余的节点插入到合并后的跳表中
	for node1 != nil {
		mergedSkipList.InsertInOrder(node1.Key, node1.DataInfo)
		node1 = node1.Next[0]
	}
	for node2 != nil {
		mergedSkipList.InsertInOrder(node2.Key, node2.DataInfo)
		node2 = node2.Next[0]
	}

	return mergedSkipList
}

// 分解跳表，确保跳表数量不超过最大限制
func (lsm *LSMTree) splitSkipListsIfNeeded() {
	// 遍历所有层级的跳表，检查是否有跳表数量超过最大限制
	for _, level := range lsm.diskLevels {
		if level.SkipListCount > lsm.maxSkipLists {
			// 如果跳表数量超过最大限制，则需要进行分解操作
			// 分解操作...
			// 这里需要根据具体的跳表实现进行分解操作
		}
	}
}
