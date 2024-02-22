package lsm

import (
	"SQL/internal/database"
	"SQL/internal/storage"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestConcurrentInsertData1(t *testing.T) {
	// 创建 LSM 树实例
	maxActiveSize := uint32(16) // 增加最大活跃内存表的大小
	maxDiskTableSize := uint32(10000)
	lsm := NewLSMTree(maxActiveSize, maxDiskTableSize)

	// 并发插入的总数量
	numInserts := 1000

	// 用于等待所有 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(numInserts)

	// 模拟并发插入操作
	for i := 0; i < numInserts; i++ {
		go func(i int) {
			defer wg.Done() // 表示当前 goroutine 已完成

			key := []byte(generateRandomKey())
			value := &DataInfo{
				DataMeta: database.DataMeta{
					Key:       key,
					Value:     []byte(fmt.Sprintf("value%d", i)),
					Extra:     []byte(fmt.Sprintf("extra%d", i)),
					KeySize:   uint32(len(key)),
					ValueSize: uint32(len(fmt.Sprintf("value%d", i))),
					ExtraSize: uint32(len(fmt.Sprintf("extra%d", i))),
					TTL:       time.Duration(rand.Intn(3600)) * time.Second,
				},
				StorageLocation: storage.StorageLocation{
					FileName: []byte("data.txt"),
					Offset:   int64(i * 100),
					Size:     100,
				},
			}
			lsm.Insert(key, value)

		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 在等待所有插入操作完成后再执行写入磁盘的操作
	lsm.SaveActiveToDiskOnExit()
	lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")

	// 等待一段时间确保数据已经写入到磁盘
	time.Sleep(100 * time.Millisecond)

	// 输出日志，查看等待的 goroutine 数量
	remaining := runtime.NumGoroutine()
	t.Logf("Remaining goroutines after test: %d", remaining)
}

func TestInsertData(t *testing.T) {
	// 创建 LSM 树实例
	maxActiveSize := uint32(10000) // 增加最大活跃内存表的大小
	maxDiskTableSize := uint32(10000)
	lsm := NewLSMTree(maxActiveSize, maxDiskTableSize)

	// 插入数据
	key := []byte("key1")
	value := &DataInfo{
		DataMeta: database.DataMeta{
			Key:       key,
			Value:     []byte("value1"),
			Extra:     []byte("extra1"),
			KeySize:   uint32(len(key)),
			ValueSize: uint32(len([]byte("value1"))),
			ExtraSize: uint32(len([]byte("extra1"))),
			TTL:       time.Duration(rand.Intn(3600)) * time.Second,
		},
		StorageLocation: storage.StorageLocation{
			FileName: []byte("data.txt"),
			Offset:   int64(0),
			Size:     100,
		},
	}
	lsm.Insert(key, value)

	// 检查插入后活跃内存表是否包含预期的数量
	expectedSize := 1
	if lsm.activeMemTable.Size != uint32(expectedSize) {
		t.Errorf("Active memory table size incorrect: expected %d, got %d", expectedSize, lsm.activeMemTable.Size)
	}
	lsm.SaveActiveToDiskOnExit()
	lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	// 等待一段时间确保数据已经写入到磁盘
	time.Sleep(100 * time.Millisecond)

	// 检查磁盘上的数据是否包含插入的数据
	// 可以根据需要编写相应的检查逻辑
}
func TestConcurrentInsertData(t *testing.T) {
	// 创建 LSM 树实例
	maxActiveSize := uint32(16) // 增加最大活跃内存表的大小
	maxDiskTableSize := uint32(10000)
	lsm := NewLSMTree(maxActiveSize, maxDiskTableSize)

	// 并发插入的总数量
	numInserts := 1000

	// 用于等待所有 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(numInserts)

	// 模拟并发插入操作
	for i := 0; i < numInserts; i++ {
		go func(i int) {
			defer wg.Done() // 表示当前 goroutine 已完成

			key := []byte(generateRandomKey())
			value := &DataInfo{
				DataMeta: database.DataMeta{
					Key:       key,
					Value:     []byte(fmt.Sprintf("value%d", i)),
					Extra:     []byte(fmt.Sprintf("extra%d", i)),
					KeySize:   uint32(len(key)),
					ValueSize: uint32(len(fmt.Sprintf("value%d", i))),
					ExtraSize: uint32(len(fmt.Sprintf("extra%d", i))),
					TTL:       time.Duration(rand.Intn(3600)) * time.Second,
				},
				StorageLocation: storage.StorageLocation{
					FileName: []byte("data.txt"),
					Offset:   int64(i * 100),
					Size:     100,
				},
			}
			lsm.Insert(key, value)

		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 在等待所有插入操作完成后再执行写入磁盘的操作

	defer lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	defer lsm.SaveActiveToDiskOnExit()
	// 等待一段时间确保数据已经写入到磁盘
	time.Sleep(100 * time.Millisecond)

	// 检查磁盘上的数据是否包含插入的数据
	// 可以根据需要编写相应的检查逻辑
}
func TestInsertSingleData(t *testing.T) {
	// 创建 LSM 树实例
	maxActiveSize := uint32(16)
	maxDiskTableSize := uint32(10000)
	lsm := NewLSMTree(maxActiveSize, maxDiskTableSize)
	dataNum := 1000
	// 创建要插入的数据
	for i := 0; i < dataNum; i++ {
		key := []byte(generateRandomKey())
		value := &DataInfo{
			DataMeta: database.DataMeta{
				Key:       key,
				Value:     []byte("value"),
				Extra:     []byte("extra"),
				KeySize:   uint32(len(key)),
				ValueSize: uint32(len("value")),
				ExtraSize: uint32(len("extra")),
				TTL:       time.Duration(60) * time.Second,
			},
			StorageLocation: storage.StorageLocation{
				FileName: []byte("data.txt"),
				Offset:   int64(0),
				Size:     100,
			},
		}
		// 插入数据到 LSM 树
		lsm.Insert(key, value)
	}

	//// 在等待所有插入操作完成后再执行写入磁盘的操作
	defer lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")

	defer lsm.SaveActiveToDiskOnExit()
	lsm.Close()
	// 等待一段时间确保数据已经写入到磁盘
	time.Sleep(100 * time.Millisecond)

}

func TestLoadFromDisk(t *testing.T) {
	// 创建一个 LSM 树实例
	lsmTree := NewLSMTree(10, 10000)

	// 定义模拟数据文件路径
	dataFilePath := "../../data/testdata/lsm_tree/test1.txt"

	// 加载模拟的数据文件到 LSM 树中
	err := lsmTree.LoadDataFromFile(dataFilePath)
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	defer lsmTree.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	defer lsmTree.SaveActiveToDiskOnExit()

}
