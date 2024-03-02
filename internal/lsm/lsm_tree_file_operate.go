package lsm

import (
	"SQL/internal/model"
	"SQL/internal/storage"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

//func (lsm *LSMTree) PrintDiskDataToFile1(filePath string) error {
//	lsm.mu.RLock()
//	defer lsm.mu.RUnlock()
//
//	// 打开文件准备写入
//	file, err := os.Create(filePath)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	// 使用 bufio.Writer 提高写入性能
//	writer := bufio.NewWriter(file)
//	defer writer.Flush()
//
//	// 遍历每个层级的跳表进行查找
//	for levelIndex, level := range lsm.diskLevels {
//		for skipListIndex, skipList := range level.SkipLists {
//			if skipList != nil {
//				writer.WriteString(fmt.Sprintf("Level %d, SkipList %d:\n", levelIndex, skipListIndex))
//				// 遍历跳表中的所有键值对并写入文件
//				skipList.ForEach(func(key []byte, value *DataInfo) bool {
//					size := strconv.FormatInt(value.Size, 10)
//					offset := strconv.FormatInt(value.Offset, 10)
//					//line := fmt.Sprintf("Key: %s, Value: %s, Extra: %s, TTL: %s, FileName: %s, Offset: %s, Size: %s\n", string(key), string(value.Value), string(value.Extra), value.TTL.XCDB_String(), string(value.FileName), offset, size)
//					line := fmt.Sprintf("Key: %s, Value: %v, Extra: %s, TTL: %s, FileName: %s, Offset: %s, Size: %s\n", string(value.Key), value.Value, string(value.Extra), value.TTL.XCDB_String(), string(value.FileName), offset, size)
//
//					_, _ = writer.WriteString(line)
//					//writer.WriteString(line)
//					return true
//				})
//
//				// 打印跳表的最大和最小键
//				writer.WriteString(fmt.Sprintf("SkipList %d, MaxKey: %s, MinKey: %s\n", skipListIndex, string(skipList.SkipListInfo.MaxKey), string(skipList.SkipListInfo.MinKey)))
//			}
//		}
//
//		// 打印当前层级的最大和最小键
//		writer.WriteString(fmt.Sprintf("InfoLevel %d, MaxKey: %s, MinKey: %s\n", levelIndex, string(level.LevelMaxKey), string(level.LevelMinKey)))
//	}
//
//	return nil
//}
//
//// LoadDataFromFile 从文件加载数据到 LSM 树中
//func (lsm *LSMTree) LoadDataFromFile1(filePath string) error {
//
//	// 打开文件，如果文件不存在则创建
//	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
//	if err != nil {
//		fmt.Println("Error opening file:", err)
//		return err
//	}
//	defer file.Close()
//	// 创建 Scanner 实例以逐行读取文件内容
//
//	//scanner := bufio.NewScanner(file)
//	scanner := bufio.NewScanner(file)
//	// 设置缓冲区大小为 64 * 1024 字节
//	const maxScanTokenSize = 1024 * 1024
//	buf := make([]byte, maxScanTokenSize)
//	scanner.Buffer(buf, maxScanTokenSize)
//
//	for scanner.Scan() {
//		line := scanner.Text()
//		if strings.HasPrefix(line, "Key") {
//			// 解析键值对信息
//			keyValue := strings.Split(line, ", ")
//			if len(keyValue) != 7 {
//				return fmt.Errorf("invalid data format: %s", line)
//			}
//
//			key := []byte(strings.Split(keyValue[0], ": ")[1])
//			value := []byte(strings.Split(keyValue[1], ": ")[1])
//			extra := []byte(strings.Split(keyValue[2], ": ")[1])
//			ttlStr := strings.Split(keyValue[3], ": ")[1]
//			ttl, err := time.ParseDuration(ttlStr)
//			if err != nil {
//				return fmt.Errorf("failed to parse TTL: %v", err)
//			}
//			fileName := []byte(strings.Split(keyValue[4], ": ")[1]) // 添加了文件名提取
//			offset, err := strconv.ParseInt(strings.Split(keyValue[5], ": ")[1], 10, 64)
//			if err != nil {
//				return fmt.Errorf("failed to parse Offset: %v", err)
//			}
//			size, err := strconv.ParseInt(strings.Split(keyValue[6], ": ")[1], 10, 64)
//			if err != nil {
//				return fmt.Errorf("failed to parse Size: %v", err)
//			}
//
//			// 创建 DataInfo 对象
//			data := DataInfo{
//				DataMeta: model.DataMeta{
//					Key:       key,
//					Value:     value,
//					Extra:     extra,
//					KeySize:   uint32(len(key)),
//					ValueSize: uint32(len(value)),
//					ExtraSize: uint32(len(extra)),
//					TTL:       ttl,
//				},
//				StorageLocation: storage.StorageLocation{
//					FileName: fileName,
//					Offset:   offset,
//					Size:     size,
//				},
//			}
//			err = lsm.Insert(key, &data)
//			if err != nil {
//				return err
//			}
//		}
//		if strings.HasPrefix(line, "Level") {
//			// 处理层级信息
//
//			// 这里可以根据需要解析和处理层级信息
//			continue
//		}
//		if strings.HasPrefix(line, "SkipList") {
//			// 处理跳表信息
//			// 这里可以根据需要解析和处理跳表信息
//			continue
//		}
//		if strings.HasPrefix(line, "InfoLevel") {
//			// 处理层级信息
//			// 这里可以根据需要解析和处理层级信息
//			continue
//		}
//
//	}
//
//	if err := scanner.Err(); err != nil {
//		return err
//	}
//
//	return nil
//}

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

// 压缩数据并写入文件
func CompressAndWriteToFile(data []byte, filePath string) error {
	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建 gzip writer
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	// 使用 bufio.Writer 提高写入性能
	writer := bufio.NewWriter(gzipWriter)
	defer writer.Flush()

	// 将数据写入压缩流
	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// 从文件读取数据并解压缩
func ReadAndDecompressFromFile(filePath string) ([]byte, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建 gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// 读取解压缩的数据
	return ioutil.ReadAll(gzipReader)
}

func (lsm *LSMTree) PrintDiskDataToFile1(filePath string) error {
	lsm.mu.RLock()
	defer lsm.mu.RUnlock()

	// 创建字符串构建器
	var sb strings.Builder

	// 遍历每个层级的跳表进行查找
	for _, level := range lsm.diskLevels {
		for _, skipList := range level.SkipLists {
			if skipList != nil {
				// 遍历跳表中的所有键值对并写入字符串构建器
				skipList.ForEach(func(key []byte, value *DataInfo) bool {
					size := strconv.FormatInt(value.Size, 10)
					offset := strconv.FormatInt(value.Offset, 10)
					line := fmt.Sprintf("Key: %s, Extra: %s, TTL: %s, FileName: %s, Offset: %s, Size: %s\n", string(value.Key), string(value.Extra), value.TTL.String(), string(value.FileName), offset, size)
					sb.WriteString(line)
					return true
				})
			}
		}
	}

	// 将数据转换为字节数组
	data := []byte(sb.String())

	// 压缩数据并写入文件
	err := CompressAndWriteToFile(data, filePath)
	if err != nil {
		return err
	}

	return nil
}

// LoadDataFromFile 从文件加载数据到 LSM 树中
func (lsm *LSMTree) LoadDataFromFile1(filePath string) error {
	// 从文件读取并解压缩数据
	compressedData, err := ReadAndDecompressFromFile(filePath)
	if err != nil {
		return err
	}

	// 将解压缩后的数据转换为字符串
	dataString := string(compressedData)

	// 使用字符串扫描器逐行解析数据
	scanner := bufio.NewScanner(strings.NewReader(dataString))
	for scanner.Scan() {
		line := scanner.Text()

		// 解析行数据
		keyValuePairs := strings.Split(line, ", ")
		if len(keyValuePairs) != 6 {
			return errors.New("invalid data format")
		}

		var key, extra, ttl, fileName, offsetStr, sizeStr string

		for _, pair := range keyValuePairs {
			splitPair := strings.Split(pair, ": ")
			if len(splitPair) != 2 {
				return errors.New("invalid data format")
			}
			switch splitPair[0] {
			case "Key":
				key = splitPair[1]
			/*case "Value":
			value = splitPair[1]*/
			case "Extra":
				extra = splitPair[1]
			case "TTL":
				ttl = splitPair[1]
			case "FileName":
				fileName = splitPair[1]
			case "Offset":
				offsetStr = splitPair[1]
			case "Size":
				sizeStr = splitPair[1]
			default:
				return errors.New("unknown key in data")
			}
		}
		ttl1, err := time.ParseDuration(ttl)
		if err != nil {
			return fmt.Errorf("failed to parse TTL: %v", err)
		}
		// 转换字符串为整数
		offsetInt, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return err
		}
		sizeInt, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			return err
		}

		// 创建 DataInfo 对象
		data := DataInfo{
			DataMeta: model.DataMeta{
				Key: []byte(key),
				//Value:     []byte(value),
				Extra:   []byte(extra),
				KeySize: uint32(len(key)),
				//	ValueSize: uint32(len(value)),
				ExtraSize: uint32(len(extra)),
				TTL:       ttl1,
			},
			StorageLocation: storage.StorageLocation{
				FileName: []byte(fileName),
				Offset:   offsetInt,
				Size:     sizeInt,
			},
		}
		// 打印解析后的数据
		//fmt.Printf("Parsed data: key=%s, value=%s, extra=%s, ttl=%s, fileName=%s, offset=%d, size=%d\n", key, value, extra, ttl, fileName, offsetInt, sizeInt)

		// 将数据信息添加到 LSM 树中

		lsm.Insert([]byte(key), &data)

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
func (lsm *LSMTree) PrintDiskDataToFile(filePath string) error {
	lsm.mu.RLock()
	defer lsm.mu.RUnlock()

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 遍历每个层级的跳表进行查找
	for _, level := range lsm.diskLevels {
		for _, skipList := range level.SkipLists {
			if skipList != nil {
				// 遍历跳表中的所有键值对并写入文件
				skipList.ForEach(func(key []byte, value *DataInfo) bool {
					size := strconv.FormatInt(value.Size, 10)
					offset := strconv.FormatInt(value.Offset, 10)
					line := fmt.Sprintf("Key: %s, Extra: %s, TTL: %s, FileName: %s, Offset: %s, Size: %s\n", string(value.Key), string(value.Extra), value.TTL.String(), string(value.FileName), offset, size)
					_, err := file.WriteString(line)
					if err != nil {
						return false // stop iteration if error occurs
					}
					return true
				})
			}
		}
	}

	return nil
}

// LoadDataFromFile 从文件加载数据到 LSM 树中
func (lsm *LSMTree) LoadDataFromFile(filePath string) error {
	// 从文件读取数据
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 将数据转换为字符串
	dataString := string(data)

	// 使用字符串扫描器逐行解析数据
	scanner := bufio.NewScanner(strings.NewReader(dataString))
	for scanner.Scan() {
		line := scanner.Text()

		// 解析行数据
		keyValuePairs := strings.Split(line, ", ")
		if len(keyValuePairs) != 6 {
			return errors.New("invalid data format")
		}

		var key, extra, ttl, fileName, offsetStr, sizeStr string

		for _, pair := range keyValuePairs {
			splitPair := strings.Split(pair, ": ")
			if len(splitPair) != 2 {
				return errors.New("invalid data format")
			}
			switch splitPair[0] {
			case "Key":
				key = splitPair[1]
			/*case "Value":
			value = splitPair[1]*/
			case "Extra":
				extra = splitPair[1]
			case "TTL":
				ttl = splitPair[1]
			case "FileName":
				fileName = splitPair[1]
			case "Offset":
				offsetStr = splitPair[1]
			case "Size":
				sizeStr = splitPair[1]
			default:
				return errors.New("unknown key in data")
			}
		}
		ttl1, err := time.ParseDuration(ttl)
		if err != nil {
			return fmt.Errorf("failed to parse TTL: %v", err)
		}
		// 转换字符串为整数
		offsetInt, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return err
		}
		sizeInt, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			return err
		}

		// 创建 DataInfo 对象
		data := DataInfo{
			DataMeta: model.DataMeta{
				Key: []byte(key),
				//Value:     []byte(value),
				Extra:   []byte(extra),
				KeySize: uint32(len(key)),
				//	ValueSize: uint32(len(value)),
				ExtraSize: uint32(len(extra)),
				TTL:       ttl1,
			},
			StorageLocation: storage.StorageLocation{
				FileName: []byte(fileName),
				Offset:   offsetInt,
				Size:     sizeInt,
			},
		}
		// 打印解析后的数据
		//fmt.Printf("Parsed data: key=%s, value=%s, extra=%s, ttl=%s, fileName=%s, offset=%d, size=%d\n", key, value, extra, ttl, fileName, offsetInt, sizeInt)

		// 将数据信息添加到 LSM 树中

		lsm.Insert([]byte(key), &data)

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
