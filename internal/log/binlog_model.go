package log

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// BinlogFile 表示一个 binlog 文件
type BinlogFile struct {
	FilePath       string   // 文件路径
	CurrFile       *os.File // 当前文件句柄
	BinLogInfo     *os.File // 这个文件记录并binlog的一些主要信息
	FileSizeMax    uint64   // 文件最大大小
	FileCurrSize   uint64   // 文件当前大小
	FileCurrNumber uint64   // 当前文件数量
	RetainDays     int      // 日志文件保留天数
}

// NewBinlogFile 初始化 BinlogFile
func NewBinlogFile(name string, fileSizeMax uint64, retainDays int) (*BinlogFile, error) {

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// 创建 BinlogFile 实例
	binlogFile := &BinlogFile{
		FilePath:    path + "/../../data/testdata/manager/" + name + "/logs/binlog/",
		FileSizeMax: fileSizeMax,
		RetainDays:  retainDays,
	}

	// 打开记录 binlog 信息的文件
	binLogInfoPath := binlogFile.FilePath + "bin_info.log"
	fmt.Println(binLogInfoPath)

	// 确保目录存在，如果不存在则递归创建
	err = os.MkdirAll(filepath.Dir(binLogInfoPath), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory for binlog info file: %v", err)
	}

	// 打开文件
	binLogInfo, err := os.OpenFile(binLogInfoPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open binlog info file: %v", err)
	}

	// 读取信息文件中的信息
	file, size, currNum, err := binlogFile.ReadInfoFromBinlogInfo()
	if err != nil {
		fmt.Println("sfd", err)
	}

	if file == nil {
		currNum = 0
		binlogFile.FileCurrNumber = currNum
		file, err = binlogFile.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open binlog info file: %v", err)
		}
		size = 0

	}

	// 设置 BinlogFile 结构体中的相应字段
	binlogFile.CurrFile = file
	binlogFile.FileCurrSize = size
	binlogFile.BinLogInfo = binLogInfo
	binlogFile.FileCurrNumber = currNum
	return binlogFile, nil
}

// ReadInfoFromBinlogInfo reads information from bin_info.log file.
func (bf *BinlogFile) ReadInfoFromBinlogInfo() (*os.File, uint64, uint64, error) {
	// Seek to the beginning of the file
	_, err := bf.BinLogInfo.Seek(0, io.SeekStart)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to seek binlog info file: %v", err)
	}

	// Use bufio.Scanner for convenient file reading
	scanner := bufio.NewScanner(bf.BinLogInfo)

	var name string
	var size uint64
	var fileNum uint64
	found := false

	// Iterate through each line in the file
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains the required information
		if strings.Contains(line, "CurrFile") && strings.Contains(line, "FileCurrSize") && strings.Contains(line, "FileCurrNumber") {
			// Extract file name and size from the line
			parts := strings.Split(line, ", ")
			for _, part := range parts {
				if strings.Contains(part, "CurrFile") {
					name = strings.TrimPrefix(part, "CurrFile: ")
				} else if strings.Contains(part, "FileCurrSize") {
					fmt.Sscanf(strings.TrimPrefix(part, "FileCurrSize: "), "%d", &size)
				} else if strings.Contains(part, "FileCurrNumber") {
					fmt.Sscanf(strings.TrimPrefix(part, "FileCurrNumber: "), "%d", &fileNum)
				}
			}
			found = true
			break
		}
	}

	// Check if the required information is found
	if !found {
		return nil, 0, 0, errors.New("required information not found in binlog info file")
	}

	// Open the file with the obtained file name
	file, err := os.Open(name)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to open current file: %v", err)
	}

	return file, size, fileNum, nil
}

// WriteInfoToBinlogInfo 将bin.log的信息文件记录到bin_info.log里面，在退出的时候进行记录
func (bf *BinlogFile) WriteInfoToBinlogInfo() error {
	// 检查当前文件是否为空，如果为空则不记录任何信息
	if bf.CurrFile == nil {
		return nil
	}

	// 获取当前文件的信息
	name := bf.CurrFile.Name()
	current, err := bf.CurrFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get current file info: %v", err)
	}
	size := current.Size()

	// 构建记录信息的字符串
	infoString := fmt.Sprintf("CurrFile: %s, FileCurrSize: %d, FileCurrNumber: %d\n", name, size, bf.FileCurrNumber)

	// 打开bin_info.log文件，如果不存在则创建
	f, err := os.OpenFile(bf.FilePath+"bin_info.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open bin_info.log: %v", err)
	}
	defer f.Close()

	// 写入信息到bin_info.log文件
	_, err = f.WriteString(infoString)
	if err != nil {
		return fmt.Errorf("failed to write info to bin_info.log: %v", err)
	}
	bf.CurrFile.Close()
	bf.BinLogInfo.Close()
	return nil
}

// Open 打开 binlog 文件
func (bf *BinlogFile) Open() (*os.File, error) {
	num := strconv.Itoa(int(bf.FileCurrNumber))
	file, err := os.OpenFile(bf.FilePath+"bin"+num+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Close 关闭 binlog 文件
func (bf *BinlogFile) Close() error {
	if bf.CurrFile != nil {
		return bf.CurrFile.Close()
	}
	return nil
}
