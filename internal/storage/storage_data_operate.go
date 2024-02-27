package storage

import (
	"SQL/internal/model"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// compressData 压缩数据
func (sm *StorageManager) compressData(data model.KeyValue) ([]byte, error) {
	// 使用 gob 包将结构体编码为字节切片
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)

	// 写入数据
	_, err := gzipWriter.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	// 关闭压缩器
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	return compressedData.Bytes(), nil
}

func DecompressData(fileName string, offset, size int64) ([]byte, error) {
	// 打开文件
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// 设置读取范围
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	// 如果指定了大小，则计算实际读取的大小
	if size > 0 && offset+size <= fileSize {
		fileSize = size
	}
	// 创建gzip.Reader
	reader, err := gzip.NewReader(io.LimitReader(file, fileSize))
	if err != nil {
		// 记录错误日志
		log.Println("Error creating gzip reader:", err)
		return nil, err
	}
	defer reader.Close()

	// 读取解压后的数据
	decompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		// 记录错误日志
		log.Println("Error reading decompressed data:", err)
		return nil, err
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
		return nil, err
	}

	return &keyValue, nil
}
