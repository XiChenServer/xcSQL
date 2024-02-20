package storage

import (
	"bytes"
	"fmt"
	"testing"
)

// TestDecompressData 测试解压数据函数
func TestDecompressData(t *testing.T) {
	fileName := "/home/zwm/GolandProjects/SQL/data/testdata/data_1.gz" // 你的存储位置文件名
	offset := int64(0)                                                 // 偏移量
	size := int64(0)                                                   // 数据大小

	// 解压数据
	decompressedData, err := DecompressData(fileName, offset, size)
	if err != nil {
		t.Fatalf("failed to decompress data: %v", err)
	}

	// 打印解压后的数据
	fmt.Println("Decompressed Data:", decompressedData)

	// 验证解压后的数据是否正确
	expectedData := generateTestData(int(size)) // 生成预期数据
	if !bytes.Equal(decompressedData, expectedData) {
		t.Errorf("decompressed data does not match expected data")
	}
}
