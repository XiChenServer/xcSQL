package log

import (
	"testing"
)

func TestBinlog(t *testing.T) {
	//// 创建测试用的文件夹和文件
	//err := os.MkdirAll("/path/to/testdata/manager/test/logs/binlog/", os.ModePerm)
	//if err != nil {
	//	t.Fatalf("failed to create test directory: %v", err)
	//}

	// 创建测试用的 BinlogFile 实例
	bf, err := NewBinlogFile("test")
	if err != nil {
		t.Fatalf("failed to create BinlogFile: %v", err)
	}

	// 测试 WriteInfoToBinlogInfo 方法
	err = bf.WriteInfoToBinlogInfo()
	if err != nil {
		t.Fatalf("WriteInfoToBinlogInfo failed: %v", err)
	}

	// 关闭测试用的 BinlogFile 实例中的文件句柄
	if bf.BinLogInfo != nil {
		bf.BinLogInfo.Close()
	}

	//// 删除测试用的文件夹和文件
	//err = os.RemoveAll("/path/to/testdata/")
	//if err != nil {
	//	t.Fatalf("failed to remove test directory: %v", err)
	//}

}
