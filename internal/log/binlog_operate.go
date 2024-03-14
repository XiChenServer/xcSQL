package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// BinlogEntry 表示哈希表的 binlog 记录条目
type BinlogEntry struct {
	Timestamp time.Time // 时间戳
	Operation string    // 操作类型（hset、hdel等）
	Key       string    // 键
	Field     []string  // 字段（仅用于某些操作，如hset）
	Value     []string  // 值（仅用于某些操作，如hset）
	TTL       uint64    // TTL（仅用于某些操作，如hset）
}

// Rotate 切换到下一个 binlog 文件
func (bf *BinlogFile) Rotate() error {
	if bf.CurrFile == nil {
		return fmt.Errorf("current file is not open")
	}

	if bf.FileCurrSize >= bf.FileSizeMax {
		if err := bf.Close(); err != nil {
			return err
		}
		num := bf.FileCurrNumber
		num++
		// 生成下一个文件名
		nextFileName := fmt.Sprintf("%s%d", bf.FilePath+"bin", num)
		bf.FileCurrNumber = num
		// 打开下一个文件
		file, err := os.OpenFile(nextFileName+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		bf.CurrFile = file
		bf.FileCurrSize = 0
	}

	return nil
}

// WriteHashEntry 向 binlog 文件中写入哈希表记录
func (bf *BinlogFile) WriteEntry(entry BinlogEntry) error {

	if bf.CurrFile == nil {
		return fmt.Errorf("current file is not open")
	}

	// 检查文件大小
	if bf.FileCurrSize >= bf.FileSizeMax {
		fmt.Println(bf.FileCurrSize, bf.FileSizeMax)
		if err := bf.Rotate(); err != nil {
			return err
		}
	}

	// 格式化记录条目
	log := entry.Format() + "\n"
	fmt.Println(len(log))
	// 写入到文件
	if _, err := bf.CurrFile.WriteString(log); err != nil {
		return err
	}
	bf.FileCurrSize += uint64(len(log))

	return nil
}

// removeOldFiles 删除超过保留天数的旧文件
func (bf *BinlogFile) removeOldFiles() error {
	files, err := filepath.Glob(fmt.Sprintf("%s.*", bf.FilePath))
	if err != nil {
		return err
	}

	// 按照文件修改时间排序
	fileInfos := make([]struct {
		name string
		time time.Time
	}, len(files))
	for i, name := range files {
		fileInfo, err := os.Stat(name)
		if err != nil {
			return err
		}
		fileInfos[i] = struct {
			name string
			time time.Time
		}{name: name, time: fileInfo.ModTime()}
	}
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].time.Before(fileInfos[j].time)
	})

	// 删除超过保留天数的文件
	for _, fileInfo := range fileInfos {
		if time.Since(fileInfo.time).Hours()/24 > float64(bf.RetainDays) {
			if err := os.Remove(fileInfo.name); err != nil {
				return err
			}
		}
	}

	return nil
}

// Format 格式化哈希表记录条目
func (be *BinlogEntry) Format() string {
	fields := strings.Join(be.Field, ",")
	values := strings.Join(be.Value, ",")
	return fmt.Sprintf("%s [%s] Key:%s Fields:%s Values:%s TTL:%s",
		be.Timestamp.Format(time.RFC3339), string(be.Operation), string(be.Key), fields, values, strconv.FormatUint(be.TTL, 10))
}
