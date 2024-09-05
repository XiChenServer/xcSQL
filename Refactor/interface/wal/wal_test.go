package wal

import (
	"strconv"
	"testing"
	"time"
)

// 整体可以使用这个命令去监控文件的大小变化
// watch -n 0.5 stat -c %s /home/zwm/GolandProjects/xcDB/Refactor/interface/wal/test.txt
func TestAddCommandToWallFileEvery(t *testing.T) {
	wal, err := CreateWalFile("test.txt", 0)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 30; i++ {
		err = wal.AddCommandToFile("add " + strconv.Itoa(i) + " to walFile")
		if err != nil {
			t.Error(err)
		}
	}
}

func TestAddCommandToWallFileNever(t *testing.T) {
	wal, err := CreateWalFile("test.txt", Never)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 10; i++ {
		err = wal.AddCommandToFile("add " + strconv.Itoa(i) + " to walFile")
		if err != nil {
			t.Error(err)
		}
		time.Sleep(time.Second)
	}
	err = wal.Close()
	if err != nil {
		return
	}

}

func TestAddCommandToWallFileSecond(t *testing.T) {
	wal, err := CreateWalFile("test.txt", Every)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 100; i++ {
		err = wal.AddCommandToFile("add " + strconv.Itoa(i) + " to walFile")
		if err != nil {
			t.Error(err)
		}
		time.Sleep(time.Second)
	}

}

// 在程序结束之后进行关闭和删除
func TestDeleteWalFile(t *testing.T) {
	wal, err := CreateWalFile("test.txt", 0)
	if err != nil {
		t.Error(err)
	}
	err = wal.Close()
	if err != nil {
		return
	}
	err = DeleteWalFile("test.txt")
	if err != nil {
		return
	}
}

//读取文件
