package database

import (
	"SQL/internal/model"
	"SQL/internal/storage"
	"SQL/logs"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// generateRandomData 生成指定长度的随机字节切片
func generateRandomData(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(rand.Intn(256)) // 生成0到255之间的随机字节
	}
	return data
}

// generateTestData 生成测试数据
func generateTestData(size int) []model.KeyValue {
	data := make([]model.KeyValue, size)
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
		data = append(data, one)
	}

	return data
}

// 简单的测试数据可以存入
func TestDB_Set(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()
	dataFilePath := "../../data/testdata/lsm_tree/test1.txt"

	// 加载模拟的数据文件到 LSM 树中
	err := db.lsm.LoadDataFromFile(dataFilePath)
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	key := []byte(generateRandomKey())
	value := []byte(generateRandomKey())
	err = db.Set(key, value)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Insert ok")
	defer db.lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	defer db.lsm.SaveActiveToDiskOnExit()
}

// 简单的测试数据可以通过解压获取到
func TestDB_GetEasy(t *testing.T) {

	fileName := "data_0.gz" // 你的存储位置文件名
	offset := int64(222)    // 偏移量
	size := int64(227)      // 数据大小
	storage := &storage.StorageManager{}
	// 解压数据
	decompressedData, err := storage.DecompressAndFillData("../../data/testdata/string_test/"+fileName, offset, size)
	if err != nil {
		t.Fatalf("failed to decompress data: %v", err)
	}
	// 打印解压后的数据
	fmt.Println("Decompressed Data:", string(decompressedData.DataMeta.Key))
}

// 简单的测试数据可以存入
func TestDB_Get(t *testing.T) {
	logs.InitLogger()
	db := NewXcDB()
	dataFilePath := "../../data/testdata/lsm_tree/test1.txt"

	// 加载模拟的数据文件到 LSM 树中
	err := db.lsm.LoadDataFromFile(dataFilePath)
	if err != nil {
		t.Fatalf("Error loading data from disk: %v", err)
	}
	key := []byte("21EQPzMyei")
	data, err := db.Get(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Get ok", string(data.DataMeta.Key), string(data.DataMeta.Value))
	defer db.lsm.PrintDiskDataToFile("../../data/testdata/lsm_tree/test1.txt")
	defer db.lsm.SaveActiveToDiskOnExit()
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
