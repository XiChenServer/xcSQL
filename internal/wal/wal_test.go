package wal

import (
	"bufio"
	"os"
	"testing"
)

func TestWAL_Write(t *testing.T) {
	// 创建临时测试文件
	//tmpFile := createTempFile(t)
	//defer os.Remove(tmpFile.Name())

	// 创建WAL实例
	wal, err := NewWAL("wal.log")
	if err != nil {
		t.Fatalf("Failed to create WAL: %v", err)
	}

	// 写入数据
	err = wal.Write("Write operation 1")
	if err != nil {
		t.Fatalf("Failed to write to WAL: %v", err)
	}

	// 刷新缓冲区
	err = wal.writer.Flush()
	if err != nil {
		t.Fatalf("Failed to flush WAL writer: %v", err)
	}

	// 读取文件内容，验证写入的数据
	//fileContent := readFileContent(t, tmpFile)
	//expectedContent := "Write operation 1\n"
	//if fileContent != expectedContent {
	//	t.Errorf("Unexpected file content. Got: %s, Expected: %s", fileContent, expectedContent)
	//}
}

func TestWAL_Recover(t *testing.T) {
	// 创建临时测试文件，并写入一些日志记录
	tmpFile := createTempFile(t)
	defer os.Remove(tmpFile.Name())
	writeLogRecords(t, tmpFile, []string{"Log 1", "Log 2", "Log 3"})

	// 创建WAL实例
	wal, err := NewWAL(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create WAL: %v", err)
	}

	// 模拟恢复操作
	err = wal.Recover()
	if err != nil {
		t.Fatalf("Failed to recover WAL: %v", err)
	}

	// 验证日志记录是否应用到数据库
	// 这里可以根据实际需求添加具体的验证逻辑
}

// 辅助函数：创建临时文件
func createTempFile(t *testing.T) *os.File {
	tmpFile, err := os.CreateTemp("", "wal_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tmpFile
}

// 辅助函数：读取文件内容
func readFileContent(t *testing.T, file *os.File) string {
	content, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Failed to read file content: %v", err)
	}
	return string(content)
}

// 辅助函数：写入日志记录到文件
func writeLogRecords(t *testing.T, file *os.File, logs []string) {
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, log := range logs {
		_, err := writer.WriteString(log + "\n")
		if err != nil {
			t.Fatalf("Failed to write log record: %v", err)
		}
	}
}
