package command

import (
	"SQL/internal/database"
	"SQL/internal/log"
	"errors"
	"fmt"

	"time"
)

// HandleSAddCommand 处理 SADD 命令
func HandleSAddCommand(parts []string, db *database.XcDB) error {
	if len(parts) < 3 {
		return errors.New("Usage: sadd [key] [member1] [member2] ... [memberN]")
	}

	key := []byte(parts[1])
	members := make([][]byte, len(parts)-2)
	for i := 2; i < len(parts); i++ {
		members[i-2] = []byte(parts[i])
	}

	return db.SAdd(key, members)
}

// HandleSRemCommand 处理 SREM 命令
func HandleSRemCommand(parts []string, db *database.XcDB) error {
	if len(parts) < 3 {
		return errors.New("Usage: srem [key] [member1] [member2] ... [memberN]")
	}

	key := []byte(parts[1])
	members := make([][]byte, len(parts)-2)
	for i := 2; i < len(parts); i++ {
		members[i-2] = []byte(parts[i])
	}

	return db.SRem(key, members)
}

// HandleSMembersCommand 处理 SMEMBERS 命令
func HandleSMembersCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: smembers [key]")
	}

	key := []byte(parts[1])
	members, err := db.SMembers(key)
	if err != nil {
		return err
	}

	logOperation(db, "smembers", key, nil)
	logResult("SMEMBERS result:", members)
	return nil
}

// HandleSIsMemberCommand 处理 SISMEMBER 命令
func HandleSIsMemberCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 3 {
		return errors.New("Usage: sismember [key] [member]")
	}

	key := []byte(parts[1])
	member := []byte(parts[2])

	isMember, err := db.SIsMember(key, member)
	if err != nil {
		return err
	}

	logOperation(db, "sismember", key, member)
	logResult("SISMEMBER result:", isMember)
	return nil
}

// HandleSCardCommand 处理 SCARD 命令
func HandleSCardCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: scard [key]")
	}

	key := []byte(parts[1])
	count, err := db.SCard(key)
	if err != nil {
		return err
	}

	logOperation(db, "scard", key, nil)
	logResult("SCARD result:", count)
	return nil
}

// logOperation 记录操作到 Binlog 中
func logOperation(db *database.XcDB, operation string, key []byte, member []byte) {
	bin := log.BinlogEntry{
		Timestamp: time.Now(),
		Operation: operation,
		Key:       string(key),
	}
	if member != nil {
		bin.Value = append(bin.Value, string(member))
	}
	db.BinLog.WriteEntry(bin)
}

// logResult 打印操作结果
func logResult(message string, result interface{}) {
	if result != nil {
		fmt.Println(message, result)
	}
}
