package storage

import (
	"fmt"
	"sync"
	"testing"
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
			_, err := storageManager.StoreData(data)
			if err != nil {
				t.Errorf("goroutine %d: failed to store data: %v", id, err)
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
func generateTestData(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}
	return data
}
