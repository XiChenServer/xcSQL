package storage

import (
	"SQL/internal/database"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// StorageManager 表示存储管理器
type StorageManager struct {
	StoragePath []byte        // 存储路径
	MaxFileSize uint64        // 文件最大大小
	CurrentFile *os.File      // 当前文件
	CurrentSize uint64        // 当前文件大小
	FileNumber  uint          // 当前文件编号
	FileLock    sync.Mutex    // 文件操作锁
	CompressBuf *bytes.Buffer // 压缩缓冲区
	CompressMtx sync.Mutex    // 压缩缓冲区锁
}

// StorageLocation 表示存储位置
type StorageLocation struct {
	FileName []byte // 文件名
	Offset   int64  // 偏移量
	Size     int64  // 数据大小
}

// NewStorageManager 创建一个新的存储管理器
func NewStorageManager(storagePath string, maxFileSize uint64) (*StorageManager, error) {
	// 创建存储路径
	err := os.MkdirAll(storagePath, 0755)
	if err != nil {
		return nil, err
	}
	// 获取当前文件编号
	fileNumber, err := getCurrentFileNumber(storagePath)
	if err != nil {
		return nil, err
	}
	// 打开或创建第一个文件
	file, err := os.OpenFile(filepath.Join(storagePath, "data_"+strconv.Itoa(int(fileNumber))+".gz"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// 获取当前文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	currentSize := uint64(fileInfo.Size())

	return &StorageManager{
		StoragePath: []byte(storagePath),
		MaxFileSize: maxFileSize,
		CurrentFile: file,
		CurrentSize: currentSize,
		FileNumber:  fileNumber,
		FileLock:    sync.Mutex{},
		CompressBuf: &bytes.Buffer{},
		CompressMtx: sync.Mutex{},
	}, nil
}

// getCurrentFileNumber 获取当前文件编号
func getCurrentFileNumber(storagePath string) (uint, error) {
	files, err := filepath.Glob(filepath.Join(storagePath, "disk_data_*.gz"))
	if err != nil {
		return 1, nil // 如果没有找到文件，返回1
	}
	return uint(len(files)), nil
}

// StoreData 将数据存储到指定位置
func (sm *StorageManager) StoreData(data *database.KeyValue) (StorageLocation, error) {
	// 压缩数据
	compressedData, err := compressData(*data)
	if err != nil {
		return StorageLocation{}, err
	}

	// 获取当前文件的偏移量和大小
	sm.FileLock.Lock()
	offset := int64(sm.CurrentSize)
	size := int64(len(compressedData))
	sm.FileLock.Unlock()

	// 如果当前文件大小超过最大限制，则创建新文件
	if offset+size > int64(sm.MaxFileSize) {
		sm.FileLock.Lock()
		sm.CurrentFile.Close()
		sm.FileNumber++
		fileName := filepath.Join(string(sm.StoragePath), "data_"+strconv.Itoa(int(sm.FileNumber))+".gz")
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			sm.FileLock.Unlock()
			return StorageLocation{}, err
		}
		sm.CurrentFile = file
		sm.CurrentSize = 0
		sm.FileLock.Unlock()
		offset = 0
	}

	// 写入数据
	sm.FileLock.Lock()
	defer sm.FileLock.Unlock() // 确保在函数返回之前释放锁

	// 移动文件指针到正确的位置
	_, err = sm.CurrentFile.Seek(offset, 0)
	if err != nil {
		return StorageLocation{}, err
	}
	_, err = sm.CurrentFile.Write(compressedData)
	if err != nil {
		return StorageLocation{}, err
	}

	// 更新当前文件大小
	sm.CurrentSize += uint64(size)

	return StorageLocation{
		FileName: []byte(sm.CurrentFile.Name()),
		Offset:   offset,
		Size:     size,
	}, nil
}

// compressData 压缩数据
func compressData(data database.KeyValue) ([]byte, error) {

	// 使用 gob 包将结构体编码为字节切片
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)

	// 写入数据
	_, err := gzipWriter.Write(buffer.Bytes())
	if err != nil {
		gzipWriter.Close()
		return nil, err
	}

	// 关闭压缩器
	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	return compressedData.Bytes(), nil
}

// DecompressData 解压数据
func DecompressData(fileName string, offset, size int64) ([]byte, error) {
	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 设置读取范围
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	// 创建gzip.Reader
	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// 读取解压后的数据
	decompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return decompressedData, nil
}
