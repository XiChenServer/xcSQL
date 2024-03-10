package wal

import (
	"bufio"
	"fmt"
	"os"
)

// 将命令行写入磁盘，并且保存这个wal的一些基本信息，保证后面可以获取
func (w *WAL) Write(data string) error {
	// 将数据写入缓冲写入器
	_, err := w.writer.WriteString(data + "\n")
	if err != nil {
		return err
	}

	// 刷新缓冲写入器，将数据刷新到磁盘
	if err := w.writer.Flush(); err != nil {
		return err
	}

	// 更新当前行数
	w.curLine++

	// 将当前行数写入信息文件
	if _, err = w.infoFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err = fmt.Fprintf(w.infoFile, "%d\n", w.curLine); err != nil {
		return err
	}

	return nil
}

// Recover 这里是对于命令的一些获取保证后面的时候，如果宕机，就可以进行处理
func (w *WAL) Recover() error {
	// 打开日志文件以供读取
	file, err := os.Open(w.logFile.Name())
	if err != nil {
		return err
	}
	defer file.Close()

	// 逐行读取日志记录
	scanner := bufio.NewScanner(file)
	for line := uint64(0); scanner.Scan(); line++ {
		if line < w.curLine {
			continue
		}
		log := scanner.Text()
		// 模拟将日志记录应用到数据库的操作
		fmt.Println("Applying log:", log)
	}

	// 检查扫描器是否出错
	if err = scanner.Err(); err != nil {
		return err
	}

	// 数据恢复完成后，删除日志文件
	if err = os.Remove(w.logFile.Name()); err != nil {
		return err
	}

	// 重置信息文件中的当前行数
	if err = w.infoFile.Truncate(0); err != nil {
		return err
	}
	if _, err = w.infoFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err = fmt.Fprintf(w.infoFile, "%d\n", w.curLine); err != nil {
		return err
	}

	return nil
}
