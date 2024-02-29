package lsm

import (
	"SQL/internal/model"
	"SQL/internal/storage"
	"SQL/logs"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// 多线程用于测试将文件写到lsm
func TestConcurrentInsertData(t *testing.T) {

	maxActiveSize := uint32(16) // 增加最大活跃内存表的大小
	maxDiskTableSize := uint32(10000)
	lsm := NewLSMTree(maxActiveSize, maxDiskTableSize, 1)

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
				DataMeta: model.DataMeta{
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
	// 保证最后的时候，可以先把活跃表中的数据写入，然后再把lsm中的数据写到磁盘
	defer lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	defer lsm.SaveActiveToDiskOnExit()
	//	defer storage.SaveStorageManager(db.storageManager, "../../data/testdata/lsm_tree/config.txt")

	// 检查磁盘上的数据是否包含插入的数据
	// 可以根据需要编写相应的检查逻辑
}

// 单线程用于测试将文件写到lsm
func TestInsertSingleData(t *testing.T) {
	// 创建 LSM 树实例
	storageManager, err := storage.NewStorageManager("../../data/testdata/storage", 4*1024) // 1MB 文件大小限制
	if err != nil {
		t.Fatalf("failed to create storage manager: %v", err)
	}
	maxActiveSize := uint32(16)
	maxDiskTableSize := uint32(10000)
	lsm := NewLSMTree(maxActiveSize, maxDiskTableSize, 1)
	dataNum := 100
	// 创建要插入的数据
	for i := 0; i < dataNum; i++ {
		//key := []byte(generateRandomKey())

		data := generateTestData()
		key := data.DataMeta.Key
		stroeLocal, err := storageManager.StoreData(&data)
		if err != nil {
			logs.SugarLogger.Error("string set fail:", err)
			return
		}
		datainfo := &DataInfo{
			DataMeta:        *data.DataMeta,
			StorageLocation: stroeLocal,
		}
		// 插入数据到 LSM 树
		lsm.Insert(key, datainfo)
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
	lsmTree := NewLSMTree(maxActiveSize, maxDiskTableSize, 1)

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

// 测试获取单个数据
func TestLoadOneData(t *testing.T) {

	maxActiveSize := uint32(16) // 增加最大活跃内存表的大小
	maxDiskTableSize := uint32(10000)
	lsmTree := NewLSMTree(maxActiveSize, maxDiskTableSize, 1)

	// 定义模拟数据文件路径
	dataFilePath := "../../data/testdata/lsm_tree/test1.txt"

	// 加载模拟的数据文件到 LSM 树中
	err := lsmTree.LoadDataFromFile(dataFilePath)
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	// 加载模拟的数据文件到 LSM 树中
	data, err := lsmTree.Get([]byte("21EQPzMyei"))
	if err != nil {
		t.Fatalf("Error loading data fatal: %v", err)
	}
	fmt.Println(string(data.Key), string(data.FileName), string(data.Value))
	defer lsmTree.PrintDiskDataToFile("../../data/testdata/lsm_tree/test2.txt")
	defer lsmTree.SaveActiveToDiskOnExit()

}

// generateTestData 生成测试数据
func generateTestData() model.KeyValue {
	var data model.KeyValue

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

	dataType := uint16(rand.Intn(100))   // 生成0到100之间的随机数据类型
	permission := uint16(rand.Intn(100)) // 生成0到100之间的随机权限控制信息
	//storageLocation := uint16(rand.Intn(100)) // 生成0到100之间的随机存储位置

	// 返回生成的随机 KeyValue 结构体实例
	one := model.KeyValue{
		DataMeta: &model.DataMeta{
			TTL:       ttl,
			Key:       key,
			Value:     value,
			Extra:     extra,
			KeySize:   uint32(len(key)),
			ValueSize: uint32(len(value)),
			ExtraSize: uint32(len(extra)),
		},
		Version:    version,
		CreateTime: createTime,
		UpdateTime: updateTime,
		AccessTime: accessTime,
		DataType:   dataType,
		DataMark:   permission,
	}
	data = one

	return data
}

// generateRandomData 生成指定长度的随机字节切片
func generateRandomData(size int) []byte {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLen := 10
	b := make([]byte, keyLen)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return b
}
