package lsm

import (
	"SQL/internal/model"
	"SQL/internal/storage"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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
					size := strconv.FormatInt(value.Size, 10)
					offset := strconv.FormatInt(value.Offset, 10)
					//line := fmt.Sprintf("Key: %s, Value: %s, Extra: %s, TTL: %s, FileName: %s, Offset: %s, Size: %s\n", string(key), string(value.Value), string(value.Extra), value.TTL.String(), string(value.FileName), offset, size)
					line := fmt.Sprintf("Key: %s, Value: %v, Extra: %s, TTL: %s, FileName: %s, Offset: %s, Size: %s\n", string(key), value.Value, string(value.Extra), value.TTL.String(), string(value.FileName), offset, size)

					writer.WriteString(line)
					return true
				})

				// 打印跳表的最大和最小键
				writer.WriteString(fmt.Sprintf("SkipList %d, MaxKey: %s, MinKey: %s\n", skipListIndex, string(skipList.SkipListInfo.MaxKey), string(skipList.SkipListInfo.MinKey)))
			}
		}

		// 打印当前层级的最大和最小键
		writer.WriteString(fmt.Sprintf("InfoLevel %d, MaxKey: %s, MinKey: %s\n", levelIndex, string(level.LevelMaxKey), string(level.LevelMinKey)))
	}

	return nil
}

// LoadDataFromFile 从文件加载数据到 LSM 树中
func (lsm *LSMTree) LoadDataFromFile(filePath string) error {

	// 打开文件，如果文件不存在则创建
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()
	// 创建 Scanner 实例以逐行读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Key") {
			// 解析键值对信息
			keyValue := strings.Split(line, ", ")
			if len(keyValue) != 7 {
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
			fileName := []byte(strings.Split(keyValue[4], ": ")[1]) // 添加了文件名提取
			offset, err := strconv.ParseInt(strings.Split(keyValue[5], ": ")[1], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse Offset: %v", err)
			}
			size, err := strconv.ParseInt(strings.Split(keyValue[6], ": ")[1], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse Size: %v", err)
			}

			// 创建 DataInfo 对象
			data := DataInfo{
				DataMeta: model.DataMeta{
					Key:       key,
					Value:     value,
					Extra:     extra,
					KeySize:   uint32(len(key)),
					ValueSize: uint32(len(value)),
					ExtraSize: uint32(len(extra)),
					TTL:       ttl,
				},
				StorageLocation: storage.StorageLocation{
					FileName: fileName,
					Offset:   offset,
					Size:     size,
				},
			}
			err = lsm.Insert(key, &data)
			if err != nil {
				return err
			}
		}
		if strings.HasPrefix(line, "Level") {
			// 处理层级信息

			// 这里可以根据需要解析和处理层级信息
			continue
		}
		if strings.HasPrefix(line, "SkipList") {
			// 处理跳表信息
			// 这里可以根据需要解析和处理跳表信息
			continue
		}
		if strings.HasPrefix(line, "InfoLevel") {
			// 处理层级信息
			// 这里可以根据需要解析和处理层级信息
			continue
		}

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// 将只读表存到lsm的磁盘之中
func (lsm *LSMTree) writeReadOnlyToDisk() {

	// 存储只读表到第一层
	lsm.storeReadOnlyToFirstLevel(lsm.readOnlyMemTable)

}

// 在程序退出时将活跃表保存到磁盘
func (lsm *LSMTree) SaveActiveToDiskOnExit() {
	lsm.readOnlyMemTable = NewSkipList(16) // 重新初始化只读内存表
	lsm.readOnlyMemTable = lsm.activeMemTable
	// 在程序退出时保存活跃表到磁盘
	defer lsm.writeReadOnlyToDisk()
}

// 在程序退出时将活跃表保存到磁盘
func (lsm *LSMTree) SaveActiveToDiskOnExit1() {

	lsm.readOnlyMemTable = lsm.activeMemTable
	// 在程序退出时保存活跃表到磁盘
	defer lsm.writeReadOnlyToDisk()
}
