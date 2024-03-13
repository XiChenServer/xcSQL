package wal

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type WAL struct {
	mu         sync.Mutex
	logFile    *os.File    // 日志文件
	curLine    uint64      // 当前行数
	infoFile   *os.File    // 用于存储额外信息的文件
	offset     int64       // 文件偏移量
	lastOffset int64       // 上次读取的文件偏移量
	cmdChan    chan string // 用于接收命令的管道
}

func NewWAL(logFile, infoFile string) (*WAL, error) {
	// 打开日志文件以供写入，如果文件不存在则创建
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// 打开信息文件以供读写，如果文件不存在则创建
	info, err := os.OpenFile(infoFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	// 读取信息文件中的当前行数和文件偏移量
	var curLine uint64
	var offset int64
	var lastOffset int64
	_, err = fmt.Fscanf(info, "curLine: %d\noffset: %d\nlastOffset: %d\n", &curLine, &offset, &lastOffset)
	if err != nil {
		// 如果无法读取到数据，则写入初始值
		initialOffset := 0
		initialLastOffset := 0
		if _, err := info.WriteString(fmt.Sprintf("curLine: %d\noffset: %d\nlastOffset: %d\n", curLine, initialOffset, initialLastOffset)); err != nil {
			return nil, err
		}
		offset = int64(initialOffset)
		lastOffset = int64(initialLastOffset)
	}

	// 创建命令管道
	cmdChan := make(chan string)

	// 返回 WAL 实例
	return &WAL{
		logFile:    file,
		curLine:    curLine,
		infoFile:   info,
		offset:     offset,
		lastOffset: lastOffset,
		cmdChan:    cmdChan,
	}, nil
}

// 将命令写入管道和日志文件，并更新信息文件中的当前行数和文件偏移量
func (wal *WAL) Write(data string) error {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// 将命令写入管道
	wal.cmdChan <- data

	// 将命令写入日志文件
	if _, err := fmt.Fprintf(wal.logFile, "%s\n", data); err != nil {
		return err
	}
	// 获取当前文件偏移量
	fileInfo, err := wal.logFile.Stat()
	if err != nil {
		return err
	}
	wal.offset = fileInfo.Size()
	// 更新当前行数
	wal.curLine++

	// 更新上次读取的文件偏移量
	wal.lastOffset = wal.offset

	// 将当前行数和文件偏移量写入信息文件
	if _, err := wal.infoFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(wal.infoFile, "curLine: %d\noffset: %d\nlastOffset: %d\n", wal.curLine, wal.offset, wal.lastOffset); err != nil {
		return err
	}

	return nil
}

// 从管道中读取下一个命令
func (wal *WAL) ReadNextCommand() (string, error) {
	// 从管道中读取命令
	cmd := <-wal.cmdChan
	return cmd, nil
}

// 恢复数据
func (wal *WAL) RecoverFromWAL() error {
	// 打开 WAL 日志文件以供读取
	file, err := os.Open(wal.logFile.Name())
	if err != nil {
		return err
	}
	defer file.Close()

	// 设置文件偏移量为上次读取到的位置
	if _, err := file.Seek(wal.lastOffset, 0); err != nil {
		return err
	}

	// 使用 bufio.Scanner 逐行读取日志记录
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 读取命令并发送到管道中
		cmd := scanner.Text()
		wal.cmdChan <- cmd
	}
	wal.lastOffset = wal.offset
	// 检查扫描器是否出错
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
