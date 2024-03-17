package command

import (
	"SQL/internal/database"
	"SQL/internal/log"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// HandleSetCommand 处理 set 命令
func HandleSetCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 3 && len(parts) != 4 {
		return errors.New("Usage: set [key] [value] [ttl]...")
	}

	key := parts[1]
	value := parts[2]
	ttl := uint64(0)
	if len(parts) == 4 {
		ttl, _ = strconv.ParseUint(parts[3], 10, 64)
	}

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "set",
		Key:       key,
		TTL:       ttl,
	}
	bin.Value = append(bin.Value, value)
	db.BinLog.WriteEntry(bin)

	// 执行操作
	return db.Set([]byte(key), []byte(value), ttl)
}

// HandleGetCommand 处理 get 命令
func HandleGetCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: get [key]")
	}

	key := parts[1]

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "get",
		Key:       key,
	}
	db.BinLog.WriteEntry(bin)

	// 执行操作
	value, err := db.Get([]byte(key))
	if err != nil {
		return err
	}

	fmt.Println("Value:", string(value))
	return nil
}

// HandleAppendCommand 处理 append 命令
func HandleAppendCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 3 {
		return errors.New("Usage: append [key] [value]")
	}

	key := parts[1]
	value := parts[2]

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "append",
		Key:       key,
	}
	bin.Value = append(bin.Value, value)
	db.BinLog.WriteEntry(bin)

	// 执行操作
	return db.Append([]byte(key), []byte(value))
}

// HandleStrlenCommand 处理 strlen 命令
func HandleStrlenCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: strlen [key]")
	}

	key := parts[1]

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "strlen",
		Key:       key,
	}
	db.BinLog.WriteEntry(bin)

	// 执行操作
	valueLen, err := db.Strlen([]byte(key))
	if err != nil {
		return err
	}

	fmt.Println("Value length:", valueLen)
	return nil
}
