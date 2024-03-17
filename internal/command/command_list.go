package command

import (
	"SQL/internal/database"
	"errors"
	"fmt"
	"strconv"
)

// HandleLPUSHCommand 处理 LPUSH 命令
func HandleLPUSHCommand(parts []string, db *database.XcDB) error {
	if len(parts) < 3 {
		return errors.New("Usage: lpush [key] [value1] [value2] ... [valueN]")
	}

	key := []byte(parts[1])
	values := make([][]byte, len(parts)-2)
	for i := 2; i < len(parts); i++ {
		values[i-2] = []byte(parts[i])
	}

	err := db.LPUSH(key, values)
	if err != nil {
		return err
	}

	return nil
}

// HandleLRANGECommand 处理 LRANGE 命令
func HandleLRANGECommand(parts []string, db *database.XcDB) error {
	if len(parts) != 4 {
		return errors.New("Usage: lrange [key] [left] [right]")
	}

	key := []byte(parts[1])
	left, err := strconv.Atoi(parts[2])
	if err != nil {
		return errors.New("Invalid left value")
	}
	right, err := strconv.Atoi(parts[3])
	if err != nil {
		return errors.New("Invalid right value")
	}

	data, err := db.LRANGE(key, left, right)
	if err != nil {
		return err
	}

	fmt.Println("LRANGE result:", data)
	return nil
}

// HandleLINDEXCommand 处理 LINDEX 命令
func HandleLINDEXCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 3 {
		return errors.New("Usage: lindex [key] [index]")
	}

	key := []byte(parts[1])
	index, err := strconv.Atoi(parts[2])
	if err != nil {
		return errors.New("Invalid index value")
	}

	data, err := db.LINDEX(key, index)
	if err != nil {
		return err
	}

	fmt.Println("LINDEX result:", data)
	return nil
}

// HandleLPOPCommand 处理 LPOP 命令
func HandleLPOPCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: lpop [key]")
	}

	key := []byte(parts[1])

	data, err := db.LPOP(key)
	if err != nil {
		return err
	}

	fmt.Println("LPOP result:", data)
	return nil
}

// HandleRPUSHCommand 处理 RPUSH 命令
func HandleRPUSHCommand(parts []string, db *database.XcDB) error {
	if len(parts) < 3 {
		return errors.New("Usage: rpush [key] [value1] [value2] ... [valueN]")
	}

	key := []byte(parts[1])
	values := make([][]byte, len(parts)-2)
	for i := 2; i < len(parts); i++ {
		values[i-2] = []byte(parts[i])
	}

	err := db.RPUSH(key, values)
	if err != nil {
		return err
	}

	return nil
}

// HandleRPOPCommand 处理 RPOP 命令
func HandleRPOPCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: rpop [key]")
	}

	key := []byte(parts[1])

	data, err := db.RPOP(key)
	if err != nil {
		return err
	}

	fmt.Println("RPOP result:", data)
	return nil
}

// HandleLLENCommand 处理 LLEN 命令
func HandleLLENCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: llen [key]")
	}

	key := []byte(parts[1])

	data, err := db.LLEN(key)
	if err != nil {
		return err
	}

	fmt.Println("LLEN result:", data)
	return nil
}
