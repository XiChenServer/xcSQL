package wal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// 刷新间隔常量
const (
	Every  = 0  // 每写一条命令都进行刷新
	Second = 1  // 每秒刷新一次
	Never  = -1 // 交给操作系统控制
)

// WAL 结构体定义
type WAL struct {
	file          *os.File
	writer        *bufio.Writer
	reader        *bufio.Reader
	flusherTicker *time.Ticker
	stopChan      chan struct{}
	filePath      string
	flushInterval time.Duration
	CommandChan   chan string
}

// CreateWalFile 创建 WAL 文件
func CreateWalFile(filePath string, flushInterval time.Duration) (*WAL, error) {
	// 使用这个方式打开文件：保证之前的数据不会发生丢失，在刷盘的时候如果出现问题的话，也不至于数据丢失不见
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to create WAL file: %w", err)
	}

	writer := bufio.NewWriter(file)
	reader := bufio.NewReader(file)

	wal := &WAL{
		file:          file,
		writer:        writer,
		reader:        reader,
		stopChan:      make(chan struct{}),
		filePath:      filePath,
		flushInterval: flushInterval,
		CommandChan:   make(chan string, 1000),
	}

	if flushInterval == Never {
		wal.flusherTicker = nil
	} else if flushInterval == Second {
		// 启动定时器定期刷新
		wal.flusherTicker = time.NewTicker(flushInterval)
		go wal.periodicFlush()
	}

	// 检查文件是否包含数据
	if err := wal.checkAndReadCommands(); err != nil {
		return nil, fmt.Errorf("failed to read WAL file: %w", err)
	}
	close(wal.CommandChan)
	return wal, nil
}

// AddCommandToFile 添加命令到 WAL 文件
func (w *WAL) AddCommandToFile(command string) error {
	if w.writer == nil {
		return fmt.Errorf("WAL file is not initialized")
	}
	_, err := w.writer.WriteString(command + "\n")
	if err != nil {
		return fmt.Errorf("failed to write command to WAL file: %w", err)
	}

	// 根据 flushInterval 决定是否立即刷新
	if w.flushInterval == Every {
		if err := w.Flush(); err != nil {
			return fmt.Errorf("failed to flush WAL file: %w", err)
		}
	}
	return nil
}

// Flush 刷新缓冲区
func (w *WAL) Flush() error {
	if w.writer == nil {
		return fmt.Errorf("WAL file is not initialized")
	}
	err := w.writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush WAL file buffer: %w", err)
	}
	return nil
}

// 定期刷新缓冲区
func (w *WAL) periodicFlush() {
	for {
		select {
		case <-w.flusherTicker.C:
			if err := w.Flush(); err != nil {
				fmt.Println("Error:", err)
			}
		case <-w.stopChan:
			w.flusherTicker.Stop()
			return
		}
	}
}

// 检查 WAL 文件是否包含数据并读取命令
func (w *WAL) checkAndReadCommands() error {
	fileInfo, err := os.Stat(w.filePath)
	if err != nil {
		return fmt.Errorf("failed to stat WAL file: %w", err)
	}
	if fileInfo.Size() == 0 {
		// 文件为空，跳过读取
		return nil
	}

	// 重新创建 reader 从文件开头读取
	_, err = w.file.Seek(0, 0)
	if err != nil {
		return err
	}
	w.reader = bufio.NewReader(w.file)
	return w.readCommands()
}

// 读取 WAL 文件中的所有命令
func (w *WAL) readCommands() error {
	for {
		line, err := w.reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to read from WAL file: %w", err)
		}
		line = strings.TrimSpace(line) // 去除行尾的换行符
		w.CommandChan <- line
	}
	return nil
}

// DeleteWalFile 删除 WAL 文件
func DeleteWalFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete WAL file: %w", err)
	}
	return nil
}

// Close 关闭 WAL 文件
func (w *WAL) Close() error {
	if w.writer != nil {
		if err := w.Flush(); err != nil {
			return err
		}
		err := DeleteWalFile(w.filePath)
		if err != nil {
			return err
		}
	}
	if w.file != nil {
		err := w.file.Close()
		if err != nil {
			return fmt.Errorf("failed to close WAL file: %w", err)
		}
	}
	if w.stopChan != nil {
		close(w.stopChan)
	}

	return nil
}
