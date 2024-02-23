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

// 多线程用于测试将文件写到lsm
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
			if i == 999 {
				fmt.Println("到了", i)
				return
			}
			fmt.Println(i, string(key))
			lsm.Insert(key, value)

		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 在等待所有插入操作完成后再执行写入磁盘的操作
	// 保证最后的时候，可以先把活跃表中的数据写入，然后再把lsm中的数据写到磁盘
	defer lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	defer lsm.SaveActiveToDiskOnExit()
	// 等待一段时间确保数据已经写入到磁盘
	time.Sleep(100 * time.Millisecond)

	// 检查磁盘上的数据是否包含插入的数据
	// 可以根据需要编写相应的检查逻辑
}

// 单线程用于测试将文件写到lsm
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

// 从文件中读取到lsm中，再从lsm中写到文件里面
func TestLoadFromDisk(t *testing.T) {

	maxActiveSize := uint32(16) // 增加最大活跃内存表的大小
	maxDiskTableSize := uint32(10000)
	lsmTree := NewLSMTree(maxActiveSize, maxDiskTableSize)

	// 定义模拟数据文件路径
	dataFilePath := "../../data/testdata/lsm_tree/test1.txt"

	// 加载模拟的数据文件到 LSM 树中
	err := lsmTree.LoadDataFromFile(dataFilePath)
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	defer lsmTree.PrintDiskDataToFile("../../data/testdata/lsm_tree/test2.txt")
	defer lsmTree.SaveActiveToDiskOnExit()

}
