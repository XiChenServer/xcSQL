package storage

import (
	"SQL/logs"
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 将StorageManager保存到文件中
func SaveStorageManager(storageManager *StorageManager, filePath string) error {
	fmt.Println(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// 使用 bufio.Writer 提高写入性能
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	name := storageManager.CurrentFile.Name()
	//name := storageManager.CurrentFile.Name()
	if name == "" {
		return errors.New("file name is nil")
	}
	writer.WriteString(fmt.Sprintf("StoragePath: %s, MaxFileSize: %d, CurrentFile: %s, CurrentSize: %d, FileNumber: %d\n", string(storageManager.StoragePath),
		storageManager.MaxFileSize, name, storageManager.CurrentSize, storageManager.FileNumber))
	return nil
}

// 将StorageManager信息从文件里面取出来
func LoadStorageManager(filePath string) (*StorageManager, error) {

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		index := strings.LastIndex(filePath, "/")
		if index != -1 {
			filePath = filePath[:index]
		}
		// 文件不存在，创建文件
		storageManager, err := NewStorageManager(filePath, 10*1024) // 1MB 文件大小限制
		if err != nil {
			logs.SugarLogger.Error("failed to create storage manager: %v", err)
			return nil, err
		}
		return storageManager, nil
		//file, err := os.Create(filePath)
		//if err != nil {
		//	return nil, err
		//}
		//defer file.Close()
	} else if err != nil {
		// 其他错误情况
		return nil, err
	}
	fmt.Println("1233")
	// 文件已存在或已创建，继续打开文件
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//file, err := os.Open(filePath)
	//if err != nil {
	//	if os.IsNotExist(err) {
	//		return nil, err // 其他错误情况
	//	}
	//	return nil, err // 其他错误情况
	//}
	//defer file.Close()
	scanner := bufio.NewScanner(file)
	var storageManager StorageManager

	if scanner.Scan() {
		line := scanner.Text()

		storageInfo := strings.Split(line, ", ")
		if len(storageInfo) != 5 {
			return nil, fmt.Errorf("invalid data format: %s", line)
		}
		StoragePath := []byte(strings.Split(storageInfo[0], ": ")[1])

		MaxFileSize := strings.Split(storageInfo[1], ": ")[1]
		maxFileSize, err := strconv.ParseUint(MaxFileSize, 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		CurrentFile := []byte(strings.Split(storageInfo[2], ": ")[1])

		CurrentSize := strings.Split(storageInfo[3], ": ")[1]
		currentSize, err := strconv.ParseUint(CurrentSize, 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		FileNumber := strings.Split(storageInfo[4], ": ")[1]
		fileNumber, err := strconv.ParseUint(FileNumber, 10, 0)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		file, err := os.OpenFile(filepath.Join(string(CurrentFile)), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		storageManager = StorageManager{
			StoragePath: StoragePath,
			MaxFileSize: maxFileSize,
			CurrentFile: file,
			CurrentSize: currentSize,
			FileNumber:  uint(fileNumber),
		}

	}

	return &storageManager, nil
}
