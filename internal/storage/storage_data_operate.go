package storage

import (
	"SQL/internal/model"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

// compressData 压缩数据
func (sm *StorageManager) compressData(data model.KeyValue) ([]byte, error) {
	// 使用 gob 包将结构体编码为字节切片
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode data: %v", err)
	}

	// 使用 gzip 压缩编码后的数据
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)

	if _, err := gzipWriter.Write(buffer.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to write compressed data: %v", err)
	}

	// 关闭 gzipWriter 并确保所有数据都被写入
	if err := gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %v", err)
	}

	return compressedData.Bytes(), nil
}

func DecompressData(fileName string, offset, size int64) ([]byte, error) {
	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 '%s': %v", fileName, err)
	}
	defer file.Close()

	// 创建一个 SectionReader 来读取指定范围的数据
	sectionReader := io.NewSectionReader(file, offset, size)

	// 读取原始数据
	var compressedData bytes.Buffer
	if _, err := io.Copy(&compressedData, sectionReader); err != nil {
		return nil, fmt.Errorf("读取压缩数据失败: %v", err)
	}

	// 创建 gzip.Reader
	reader, err := gzip.NewReader(bytes.NewReader(compressedData.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("创建 gzip 读取器失败: %v", err)
	}
	defer reader.Close()

	// 读取解压后的数据
	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取解压后的数据失败: %v", err)
	}

	return decompressedData, nil
}

func (sm *StorageManager) DecompressAndFillData(fileName string, offset, size int64) (*model.KeyValue, error) {
	// 解压数据
	decompressedData, err := DecompressData(fileName, offset, size)
	if err != nil {
		return nil, err
	}

	// 解码数据到 KeyValue 结构体
	var keyValue model.KeyValue
	err = gob.NewDecoder(bytes.NewReader(decompressedData)).Decode(&keyValue)
	if err != nil {
		fmt.Printf("解码数据失败: %v\n", err)
		return nil, err
	}

	return &keyValue, nil
}
