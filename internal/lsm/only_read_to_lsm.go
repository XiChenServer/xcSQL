package lsm

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
)

// 将活跃内存表转换为只读表
func (lsm *LSMTree) convertActiveToReadOnly() {

	lsm.readOnlyMemTable = lsm.activeMemTable
	lsm.activeMemTable = NewSkipList(16) // 重新初始化活跃内存表

}

// 更新一个层中键的最大和最小的问题
func (lsm *LSMTree) updateLevelMinMaxKeys(currentLevel *LevelInfo, selectedSkipList *SkipList) error {
	// 检查传入参数的有效性
	if currentLevel == nil || selectedSkipList == nil || selectedSkipList.SkipListInfo == nil {
		fmt.Println("Invalid currentLevel or selectedSkipList.")
		err := errors.New("Invalid currentLevel or selectedSkipList.")
		return err
	}

	// 获取跳表的最小键和最大键
	minKey := selectedSkipList.SkipListInfo.MinKey
	maxKey := selectedSkipList.SkipListInfo.MaxKey

	// 如果跳表为空，则直接返回
	if minKey == nil || maxKey == nil {
		err := errors.New("skiplist is nil")
		return err
	}

	// 如果当前层级的最小键为空或者跳表的最小键小于当前层级的最小键，则更新最小键
	if len(currentLevel.LevelMinKey) == 0 || bytes.Compare(minKey, currentLevel.LevelMinKey) < 0 {
		currentLevel.LevelMinKey = minKey
	}

	// 如果当前层级的最大键为空或者跳表的最大键大于当前层级的最大键，则更新最大键
	if len(currentLevel.LevelMaxKey) == 0 || bytes.Compare(maxKey, currentLevel.LevelMaxKey) > 0 {
		currentLevel.LevelMaxKey = maxKey
	}
	return nil
}

func (lsm *LSMTree) storeReadOnlyToFirstLevel(skipList *SkipList) error {
	// 如果第一层还有空间，则直接存储到第一层
	if lsm.diskLevels[0].SkipListCount >= lsm.diskLevels[0].LevelMaxSkipListCount {

		// 如果第一层已满，则随机选择一个跳表存储到下一层
		randomIndex := rand.Intn(int(lsm.diskLevels[0].SkipListCount))
		selectedSkipList := lsm.diskLevels[0].SkipLists[randomIndex]
		lsm.moveSkipListDown(0, randomIndex, selectedSkipList)

	}

	err := lsm.keepLsmLevelOrderly(0, skipList)
	if err != nil {
		return err
	}
	return nil
}

// 将跳表从一个级别移动到下一个较低级别，并删除原来的位置
func (lsm *LSMTree) moveSkipListDown(levelIndex, skipListIndex int, skipList *SkipList) error {
	// 如果下一层有空间，则存储到下一层
	if levelIndex+1 < len(lsm.diskLevels) && lsm.diskLevels[levelIndex+1].SkipListCount >= lsm.diskLevels[levelIndex+1].LevelMaxSkipListCount {
		// 如果下一层也满了，则随机选择一个跳表存储到下一层
		randomIndex := rand.Intn(int(lsm.diskLevels[levelIndex+1].SkipListCount))
		selectedSkipList := lsm.diskLevels[levelIndex+1].SkipLists[randomIndex]
		err := lsm.moveSkipListDown(levelIndex+1, randomIndex, selectedSkipList)
		if err != nil {
			return err
		}
	}
	// 如果下一层有空间，则将选定的跳表存储到下一层

	err := lsm.keepLsmLevelOrderly(levelIndex+1, skipList)
	if err != nil {
		return err
	}
	err = lsm.deleteSkipList(levelIndex, skipListIndex)
	if err != nil {
		return err
	}
	return nil
}

// 删除指定层级的跳表
func (lsm *LSMTree) deleteSkipList(levelIndex, skipListIndex int) error {
	// 检查待删除的索引是否有效
	if skipListIndex < 0 || skipListIndex >= len(lsm.diskLevels[levelIndex].SkipLists) {
		fmt.Println("Invalid skipListIndex.")
		err := errors.New("Invalid skipListIndex.")
		return err
	}

	// 将待删除的跳表从切片中移除
	skipLists := lsm.diskLevels[levelIndex].SkipLists
	copy(skipLists[skipListIndex:], skipLists[skipListIndex+1:])
	lsm.diskLevels[levelIndex].SkipLists = skipLists[:len(skipLists)-1]

	lsm.diskLevels[levelIndex].SkipListCount--
	return nil
}
