package storage

import (
	"SQL/internal/database"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestConcurrentStoreData(t *testing.T) {
	const numRoutines = 100   // 并发测试的协程数量
	const testDataSize = 1000 // 测试数据的大小

	// 创建存储管理器
	storageManager, err := NewStorageManager("../../data/testdata", 4*1024) // 1MB 文件大小限制
	if err != nil {
		t.Fatalf("failed to create storage manager: %v", err)
	}
	//defer func() {
	//	// 清理测试数据
	//	err := os.RemoveAll("testdata")
	//	if err != nil {
	//		t.Fatalf("failed to clean up test data: %v", err)
	//	}
	//}()

	// 并发写入测试数据
	var wg sync.WaitGroup
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			data := generateTestData(testDataSize)
			for _, v := range data {
				_, err := storageManager.StoreData(v)
				if err != nil {
					t.Errorf("goroutine %d: failed to store data: %v", id, err)
				}
			}

		}(i)
	}
	wg.Wait()

	// 检查存储结果
	fileInfo, err := storageManager.CurrentFile.Stat()
	if err != nil {
		t.Fatalf("failed to get file info: %v", err)
	}
	fmt.Printf("Total data stored: %d bytes\n", fileInfo.Size())
}

// generateTestData 生成测试数据
func generateTestData(size int) []database.KeyValuePair {
	data := make([]database.KeyValuePair, size)
	for i := 0; i < size; i++ {

		// generateRandomKeyValuePair 生成随机的 KeyValuePair 结构体实例
		// 生成随机的键、值和额外信息
		key := generateRandomData(10)   // 生成长度为10的随机字节切片作为键
		value := generateRandomData(20) // 生成长度为20的随机字节切片作为值
		extra := generateRandomData(5)  // 生成长度为5的随机字节切片作为额外信息

		// 生成随机的 TTL、版本号和时间
		ttl := time.Duration(rand.Intn(3600)) * time.Second // 生成0到3600秒之间的随机 TTL
		version := rand.Uint32()                            // 生成随机的版本号
		createTime := time.Now()                            // 记录当前时间作为创建时间
		updateTime := time.Now()                            // 记录当前时间作为修改时间
		accessTime := time.Now()                            // 记录当前时间作为访问时间

		// 生成随机的标签、数据类型、权限控制信息和存储位置
		tags := generateRandomData(8)             // 生成长度为8的随机字节切片作为标签
		dataType := uint16(rand.Intn(100))        // 生成0到100之间的随机数据类型
		permission := uint16(rand.Intn(100))      // 生成0到100之间的随机权限控制信息
		storageLocation := uint16(rand.Intn(100)) // 生成0到100之间的随机存储位置

		// 返回生成的随机 KeyValuePair 结构体实例
		one := database.KeyValuePair{
			KvType: &database.KvType{
				Key:       key,
				Value:     value,
				Extra:     extra,
				KeySize:   uint32(len(key)),
				ValueSize: uint32(len(value)),
				ExtraSize: uint32(len(extra)),
			},
			TTL:             ttl,
			Version:         version,
			CreateTime:      createTime,
			UpdateTime:      updateTime,
			AccessTime:      accessTime,
			Tags:            tags,
			DataType:        dataType,
			Permission:      permission,
			StorageLocation: storageLocation,
		}
		data = append(data, one)
	}

	return data
}

// generateRandomData 生成指定长度的随机字节切片
func generateRandomData(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(rand.Intn(256)) // 生成0到255之间的随机字节
	}
	return data
}
