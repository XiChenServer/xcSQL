package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_config(t *testing.T) {
	//// 打开文件
	//cwd, err := os.Getwd()
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//fmt.Println("Current working directory:", cwd)
	//file, err := os.OpenFile("/home/zwm/GolandProjects/xcDB/test/cmd/log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	//if err != nil {
	//	t.Errorf("Failed to open file: %v", err)
	//	return
	//}
	//defer file.Close()
	//
	//// 检查文件是否成功打开
	//if fileInfo, _ := file.Stat(); fileInfo == nil {
	//	t.Error("File should be opened, but it's not.")
	//}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	relativePath := filepath.Join(cwd, "your", "relative", "path")
	fmt.Println("Relative Path:", relativePath)

}
