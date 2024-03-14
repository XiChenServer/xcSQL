package lsm

import (
	"SQL/internal/model"
	"SQL/internal/storage"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

// TestDataInfoGenerationAndWrite 测试 DataInfo 生成和写入
func TestDataInfoGenerationAndWrite(t *testing.T) {
	// 并发写入的 goroutine 数量
	concurrency := 100

	// 创建跳表实例
	sl := NewSkipList(16)

	// 生成测试数据并插入跳表
	var testData []DataInfo
	for i := 0; i < concurrency; i++ {
		data := DataInfo{
			DataMeta: model.DataMeta{
				Key:       []byte(generateRandomKey()),
				Extra:     []byte(fmt.Sprintf("extra%d", i)),
				KeySize:   uint32(len(fmt.Sprintf("key%d", i))),
				ExtraSize: uint32(len(fmt.Sprintf("extra%d", i))),
				TTL:       time.Duration(rand.Intn(3600)) * time.Second, // 随机生成 TTL
			},
			StorageLocation: storage.StorageLocation{
				FileName: []byte("data.txt"),
				Offset:   int64(i * 100), // 假设每条数据占用 100 字节
				Size:     100,
			},
		}
		testData = append(testData, data)
		sl.InsertInOrder(data.Key, &data)
	}

	// 并发写入数据到文件
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for _, data := range testData {
		go func(d DataInfo) {
			defer wg.Done()
			writeDataToFile(d, sl)
		}(data)
	}
	wg.Wait()

	// 将跳表内容写入文件
	file, err := os.OpenFile("../../data/testdata/skiplist/skiplist_content.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	// 将跳表内容写入文件
	file, err = os.OpenFile("../../data/testdata/skiplist/skiplist_content.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	// 遍历跳表中的每个节点并将数据写入文件
	sl.ForEach(func(key []byte, value *DataInfo) bool {
		line := fmt.Sprintf("Key: %s,  Extra: %s, KeySize: %d, ExtraSize: %d, TTL: %s, FileName: %s, Offset: %d, Size: %d\n",
			key, value.Extra, value.KeySize, value.ExtraSize, value.TTL, value.FileName, value.Offset, value.Size)
		if _, err := file.WriteString(line); err != nil {
			fmt.Printf("failed to write to file: %v\n", err)
			return false
		}
		return true
	})
}

// writeDataToFile 将 DataInfo 写入文件
func writeDataToFile(data DataInfo, sl *SkipList) {
	// 打开文件，如果不存在则创建
	file, err := os.OpenFile(string("../../data/testdata/skiplist/test1.txt"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	// 将 DataInfo 格式化为字符串
	line := fmt.Sprintf("Key: %s, Extra: %s, KeySize: %d, ExtraSize: %d, TTL: %s, FileName: %s, Offset: %d, Size: %d\n",
		data.Key, data.Extra, data.KeySize, data.ExtraSize, data.TTL, data.FileName, data.Offset, data.Size)
	key := fmt.Sprintf("keyMax: %s, keyMin: %s", sl.SkipListInfo.MaxKey, sl.SkipListInfo.MinKey)
	// 写入数据到文件
	_, err = file.WriteString(line)
	if err != nil {
		fmt.Printf("failed to write to file: %v\n", err)
	}
	_, err = file.WriteString(key)
	if err != nil {
		fmt.Printf("failed to write to file: %v\n", err)
	}
}

// go test -run=^$ -bench=. -benchmem
func BenchmarkDataInfoGenerationAndWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestDataInfoGenerationAndWrite(nil)
	}
}

// generateRandomKey 生成随机键值
func generateRandomKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLen := 10
	b := make([]byte, keyLen)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// TestDataInfoGenerationAndWrite 测试 DataInfo 生成和写入
func TestDataInfoGenerationAndWrite1(t *testing.T) {
	// 创建跳表实例
	sl := NewSkipList(16)

	// 创建 WaitGroup 以等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 设置并发数
	concurrency := 10000

	// 设置互斥锁以保护对跳表的并发写入
	var mu sync.Mutex

	// 启动并发写入 goroutine
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			data := DataInfo{
				DataMeta: model.DataMeta{
					Key: []byte(generateRandomKey()),

					Extra:   []byte(fmt.Sprintf("extra%d", i)),
					KeySize: uint32(len(fmt.Sprintf("key%d", i))),

					ExtraSize: uint32(len(fmt.Sprintf("extra%d", i))),
					TTL:       time.Duration(rand.Intn(3600)) * time.Second, // 随机生成 TTL
				},
				StorageLocation: storage.StorageLocation{
					FileName: []byte("data.txt"),
					Offset:   int64(i * 100), // 假设每条数据占用 100 字节
					Size:     100,
				},
			}

			// 插入数据到跳表
			mu.Lock()
			defer mu.Unlock()
			sl.InsertInOrder(data.Key, &data)
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 将跳表内容写入文件
	file, err := os.OpenFile("../../data/testdata/skiplist/skiplist_content.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	// 遍历跳表中的每个节点并将数据写入文件
	sl.ForEach(func(key []byte, value *DataInfo) bool {
		line := fmt.Sprintf("Key: %s,  Extra: %s, KeySize: %d, ExtraSize: %d, TTL: %s, FileName: %s, Offset: %d, Size: %d\n",
			key, value.Extra, value.KeySize, value.ExtraSize, value.TTL, value.FileName, value.Offset, value.Size)
		if _, err := file.WriteString(line); err != nil {
			fmt.Printf("failed to write to file: %v\n", err)
			return false
		}
		return true
	})
}
