package lsm

import (
	"SQL/internal/database"
	"SQL/internal/storage"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestConcurrentInsert(t *testing.T) {
	// 并发插入的 goroutine 数量
	concurrency := 100

	// 创建 LSM 树实例
	maxActiveSize := uint32(1000)
	maxDiskTableSize := uint32(10000)
	lsm := NewLSMTree(maxActiveSize, maxDiskTableSize)

	// 模拟并发插入操作
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			numInserts := 100
			for j := 0; j < numInserts; j++ {
				key := []byte(fmt.Sprintf("key%d", j))
				data := DataInfo{
					DataMeta: database.DataMeta{
						Key:       []byte(generateRandomKey()),
						Value:     []byte(fmt.Sprintf("value%d", i)),
						Extra:     []byte(fmt.Sprintf("extra%d", i)),
						KeySize:   uint32(len(fmt.Sprintf("key%d", i))),
						ValueSize: uint32(len(fmt.Sprintf("value%d", i))),
						ExtraSize: uint32(len(fmt.Sprintf("extra%d", i))),
						TTL:       time.Duration(rand.Intn(3600)) * time.Second, // 随机生成 TTL
					},
					StorageLocation: storage.StorageLocation{
						FileName: []byte("data.txt"),
						Offset:   int64(i * 100), // 假设每条数据占用 100 字节
						Size:     100,
					},
				}
				lsm.Insert(key, &data)
			}
		}()
	}
	wg.Wait()

	// 检查插入后活跃内存表是否包含预期的数量
	expectedSize := concurrency * 100
	if lsm.activeMemTable.Size != uint32(expectedSize) {
		t.Errorf("Active memory table size incorrect: expected %d, got %d", expectedSize, lsm.activeMemTable.Size)
	}

	// 将活跃内存表转换为只读内存表并写入磁盘
	lsm.convertActiveToReadOnly()
	lsm.writeReadOnlyToDisk()

	// 检查只读内存表是否被清空
	if lsm.readOnlyMemTable != nil {
		t.Error("Read-only memory table not cleared after conversion")
	}

	// TODO: 添加其他测试步骤，例如验证数据是否正确写入磁盘等
}
func TestLSMTree0_InsertConcurrentInsert(t *testing.T) {
	// 创建 LSM 树实例
	maxActiveSize := uint32(1000)
	maxDiskTableSize := uint32(10000)
	lsm := NewLSMTree(maxActiveSize, maxDiskTableSize)

	// 模拟插入操作
	numInserts := 100 * 100 // 并发插入的总数量
	for i := 0; i < numInserts; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		data := DataInfo{
			DataMeta: database.DataMeta{
				Key:       key,
				Value:     []byte(fmt.Sprintf("value%d", i)),
				Extra:     []byte(fmt.Sprintf("extra%d", i)),
				KeySize:   uint32(len(fmt.Sprintf("key%d", i))),
				ValueSize: uint32(len(fmt.Sprintf("value%d", i))),
				ExtraSize: uint32(len(fmt.Sprintf("extra%d", i))),
				TTL:       time.Duration(rand.Intn(3600)) * time.Second, // 随机生成 TTL
			},
			StorageLocation: storage.StorageLocation{
				FileName: []byte("data.txt"),
				Offset:   int64(i * 100), // 假设每条数据占用 100 字节
				Size:     100,
			},
		}
		lsm.Insert(key, &data)
	}

	// 检查插入后活跃内存表是否包含预期的数量
	expectedSize := numInserts
	if lsm.activeMemTable.Size != uint32(expectedSize) {
		t.Errorf("Active memory table size incorrect: expected %d, got %d", expectedSize, lsm.activeMemTable.Size)
	}
}
func TestLSMTree1_InsertConcurrentInsert(t *testing.T) {
	// 创建 LSM 树实例
	maxActiveSize := uint32(10000) // 增加最大活跃内存表的大小
	maxDiskTableSize := uint32(10000)
	lsm := NewLSMTree(maxActiveSize, maxDiskTableSize)

	// 并发插入的总数量
	numInserts := 100 * 100

	// 用于等待所有 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(numInserts)

	// 模拟并发插入操作
	for i := 0; i < numInserts; i++ {
		go func(i int) {
			defer wg.Done() // 表示当前 goroutine 已完成

			key := []byte(fmt.Sprintf("key%d", i))
			data := DataInfo{
				DataMeta: database.DataMeta{
					Key:       key,
					Value:     []byte(fmt.Sprintf("value%d", i)),
					Extra:     []byte(fmt.Sprintf("extra%d", i)),
					KeySize:   uint32(len(fmt.Sprintf("key%d", i))),
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
			lsm.Insert(key, &data)
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 检查插入后活跃内存表是否包含预期的数量
	expectedSize := numInserts
	if lsm.activeMemTable.Size != uint32(expectedSize) {
		t.Errorf("Active memory table size incorrect: expected %d, got %d", expectedSize, lsm.activeMemTable.Size)
	}
	lsm.SaveActiveToDiskOnExit()
	lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
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
	maxActiveSize := uint32(10) // 增加最大活跃内存表的大小
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
	lsm.SaveActiveToDiskOnExit()
	lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	// 等待一段时间确保数据已经写入到磁盘

	// 等待一段时间确保数据已经写入到磁盘
	time.Sleep(100 * time.Millisecond)

	// 检查磁盘上的数据是否包含插入的数据
	// 可以根据需要编写相应的检查逻辑
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
	lsmTree.SaveActiveToDiskOnExit()
	lsmTree.PrintDiskDataToFile("../../data/testdata/lsm_tree/test2.txt")

}
