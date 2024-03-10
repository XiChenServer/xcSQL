package wal

import (
	"bufio"
	"fmt"
	"os"
)

type WAL struct {
	logFile  *os.File // 日志文件
	writer   *bufio.Writer
	curLine  uint64   // 当前行数
	infoFile *os.File // 用于存储额外信息的文件
}

func NewWAL(logFile, infoFile string) (*WAL, error) {
	// 打开日志文件以供写入，如果文件不存在则创建
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// 创建缓冲写入器
	writer := bufio.NewWriter(file)

	// 打开信息文件以供读写，如果文件不存在则创建
	info, err := os.OpenFile(infoFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	// 读取信息文件中的当前行数
	var curLine uint64
	_, err = fmt.Fscanf(info, "%d", &curLine)
	if err != nil {
		if _, err := info.WriteString("0\n"); err != nil {
			return nil, err
		}
		curLine = 0
	}

	return &WAL{
		logFile:  file,
		writer:   writer,
		curLine:  curLine,
		infoFile: info,
	}, nil
}
