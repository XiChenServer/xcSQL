package command

import (
	"SQL/internal/database"
	"SQL/internal/log"
	"errors"
	"fmt"
	"time"
)

// HandleHSetCommand 处理 HSet 命令
func HandleHSetCommand(parts []string, db *database.XcDB) error {
	if len(parts) < 3 {
		return errors.New("Usage: hset [key] [field1] [value1] ... [fieldN] [valueN]")
	}

	key := []byte(parts[1])
	if len(parts)%2 != 0 {
		return errors.New("Odd number of arguments for field-value pairs")
	}

	valueMap := make(map[string]string)
	for i := 2; i < len(parts); i += 2 {
		field := parts[i]
		value := parts[i+1]
		valueMap[field] = value
	}

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "hset",
		Key:       string(key),
		Field:     make([]string, 0, len(valueMap)),
		Value:     make([]string, 0, len(valueMap)),
		TTL:       0, // 根据需要修改 TTL 的值
	}

	// 将字段和值分别添加到记录中
	for field, value := range valueMap {
		bin.Field = append(bin.Field, field)
		bin.Value = append(bin.Value, value)
	}

	db.BinLog.WriteEntry(bin)

	// 执行操作
	err := db.HSet(key, valueMap)
	if err != nil {
		return err
	}

	return nil
}

// HandleHGetCommand 处理 HGet 命令
func HandleHGetCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 3 {
		return errors.New("Usage: hget [key] [field]")
	}

	key := []byte(parts[1])
	field := parts[2]

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "hget",
		Key:       string(key),
	}
	db.BinLog.WriteEntry(bin)

	// 执行操作
	value, err := db.HGet(key, field)
	if err != nil {
		return err
	}

	fmt.Println("Field value:", string(value))
	return nil
}

// HandleHGetAllCommand 处理 HGETALL 命令
func HandleHGetAllCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: hgetall [key]")
	}

	key := []byte(parts[1])

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "hgetall",
		Key:       string(key),
	}
	db.BinLog.WriteEntry(bin)

	// 执行操作
	values, err := db.HGETALL(key)
	if err != nil {
		return err
	}

	fmt.Println("All fields and values:", values)
	return nil
}

// HandleHDelCommand 处理 HDEL 命令
func HandleHDelCommand(parts []string, db *database.XcDB) error {
	if len(parts) < 3 {
		return errors.New("Usage: hdel [key] [field1] [field2] ...")
	}

	key := []byte(parts[1])
	fields := parts[2:]

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "hdel",
		Key:       string(key),
		Field:     fields,
	}
	db.BinLog.WriteEntry(bin)

	// 执行操作
	err := db.HDel(key, fields...)
	if err != nil {
		return err
	}

	return nil
}

// HandleHExistsCommand 处理 HEXISTS 命令
func HandleHExistsCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 3 {
		return errors.New("Usage: hexists [key] [field]")
	}

	key := []byte(parts[1])
	field := parts[2]

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "hexists",
		Key:       string(key),
	}
	bin.Field = append(bin.Field, field)
	db.BinLog.WriteEntry(bin)

	// 执行操作
	exists, err := db.HExists(key, field)
	if err != nil {
		return err
	}

	if exists {
		fmt.Println("Field exists in hash")
	} else {
		fmt.Println("Field does not exist in hash")
	}
	return nil
}

// HandleHKeysCommand 处理 HKEYS 命令
func HandleHKeysCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: hkeys [key]")
	}

	key := []byte(parts[1])

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "hkeys",
		Key:       string(key),
	}
	db.BinLog.WriteEntry(bin)

	// 执行操作
	keys, err := db.HKeys(key)
	if err != nil {
		return err
	}

	fmt.Println("Hash keys:", keys)
	return nil
}

// HandleHValsCommand 处理 HVALS 命令
func HandleHValsCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: hvals [key]")
	}

	key := []byte(parts[1])

	// 将合法操作记录到 Binlog 中
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: "hvals",
		Key:       string(key),
	}
	db.BinLog.WriteEntry(bin)

	// 执行操作
	vals, err := db.HVals(key)
	if err != nil {
		return err
	}

	fmt.Println("Hash values:", vals)
	return nil
}
