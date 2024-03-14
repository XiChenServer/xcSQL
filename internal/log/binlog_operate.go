package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// BinlogEntry 表示一个 binlog 记录条目
type BinlogEntry struct {
	Timestamp time.Time // 时间戳
	Operation string    // 操作类型（插入、更新、删除等）
	Key       string    // 键
	Value     string    // 值
	TTL       uint64    //额外的数据
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

// WriteEntry 向 binlog 文件中写入记录
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

//// currFileSize 获取当前文件的大小
//func (bf *BinlogFile) currFileSize() (uint64, error) {
//	fileInfo, err := bf.CurrFile.Stat()
//	if err != nil {
//		return 0, err
//	}
//	return uint64(fileInfo.Size()), nil
//}

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

// Format 格式化记录条目
func (be *BinlogEntry) Format() string {
	return fmt.Sprintf("%s [%s] Key:%s Value:%s TTL:%s",
		be.Timestamp.Format(time.RFC3339), string(be.Operation), string(be.Key), string(be.Value), strconv.FormatUint(be.TTL, 10))
}
