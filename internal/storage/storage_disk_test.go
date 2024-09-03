package storage

import (
	"SQL/internal/model"
	"fmt"
	"github.com/magiconair/properties/assert"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

func TestConcurrentStoreData(t *testing.T) {
	const numRoutines = 10    // 并发测试的协程数量
	const testDataSize = 1000 // 测试数据的大小

	// 创建存储管理器
	storageManager, err := NewStorageManager("../../data/testdata", 4*1024) // 1MB 文件大小限制
	if err != nil {
		t.Fatalf("failed to create storage manager: %v", err)
	}

	// 并发写入测试数据
	var wg sync.WaitGroup
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			data := generateTestData(testDataSize)
			for _, v := range data {
				_, err := storageManager.StoreData(&v)
				if err != nil {
					t.Errorf("goroutine %d: failed to store data: %v", id, err)
				}
			}

		}(i)
	}
	wg.Wait()

}

func TestConcurrentStoreData1(t *testing.T) {

	// 创建存储管理器
	storageManager, err := NewStorageManager("../../data/testdata1", 4*1024) // 1MB 文件大小限制
	if err != nil {
		t.Fatalf("failed to create storage manager: %v", err)
	}

	data := generateTestData(1)
	for i := 0; i < len(data); i++ {
		p, err := storageManager.StoreData(&data[i])
		fmt.Println(p.Offset, string(p.FileName), p.Size)
		if err != nil {
			t.Errorf("failed to store data: %v", err)
		}
	}

}

// TestStoreData 测试 StoreData 方法
func TestStoreData(t *testing.T) {
	// 创建存储管理器
	storageManager, err := NewStorageManager("../../data/testdata1", 4*1024) // 4KB 文件大小限制
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	// 生成测试数据
	data := generateTestData(1)
	var sm StorageManager
	for _, v := range data {
		// 存储数据
		p, err := storageManager.StoreData(&v)
		if err != nil {
			t.Errorf("StoreData failed: %v", err)
			continue
		}

		// 解压数据
		decompressedData, err := (&sm).DecompressAndFillData(string(p.FileName), int64(p.Offset), int64(p.Size))
		if err != nil {
			t.Errorf("DecompressData failed: %v", err)
			continue
		}
		fmt.Println(decompressedData)
		// 断言数据一致性
		assert.Equal(t, v, decompressedData, "Stored and decompressed data should match")

		// 打印解压后的数据用于调试
		// fmt.Printf("Decompressed data: %s\n", decompressedData)
	}

	// 清理测试文件
	os.RemoveAll("../../data/testdata1")
}

// generateTestData 生成测试数据
func generateTestData(size int) []model.KeyValue {
	data := make([]model.KeyValue, 0, size)
	for i := 0; i < size; i++ {

		// generateRandomKeyValuePair 生成随机的 KeyValue 结构体实例
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

		permission := uint16(rand.Intn(100)) // 生成0到100之间的随机权限控制信息
		//storageLocation := uint16(rand.Intn(100)) // 生成0到100之间的随机存储位置

		// 返回生成的随机 KeyValue 结构体实例
		one := model.KeyValue{
			DataMeta: &model.DataMeta{
				TTL:       ttl,
				Key:       key,
				Extra:     extra,
				KeySize:   uint32(len(key)),
				ExtraSize: uint32(len(extra)),
			},
			ValueSize:  uint32(len(value)),
			Value:      value,
			Version:    version,
			CreateTime: createTime,
			UpdateTime: updateTime,
			AccessTime: accessTime,
			DataType:   model.XCDB_String,
			DataMark:   permission,
		}
		data = append(data, one)
	}

	return data
}

// generateRandomData 生成指定长度的随机字节切片
// generateRandomData 生成指定长度的随机字节切片，包含数字 0-9 和字母 a-A
func generateRandomData(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, size)
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i := range data {
		data[i] = charset[rand.Intn(len(charset))]
	}
	return data
}
