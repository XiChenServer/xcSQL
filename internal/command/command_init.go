package command

import (
	"SQL/internal/database"
	"SQL/internal/log"
	"errors"
	"fmt"
	"strconv"
)

func HandleSetCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 3 && len(parts) != 4 {
		return errors.New("Usage: set [key] [value] [ttl]...")
	}
	key := []byte(parts[1])
	value := []byte(parts[2])
	ttl := uint64(0)
	if len(parts) == 4 {
		ttl, _ = strconv.ParseUint(parts[3], 10, 64)
	}
	var bin = log.BinlogEntry{
		Operation: []byte("set"),
		Key:       []byte(parts[1]),
		Value:     []byte(parts[2]),
		Extra:     []byte(parts[3]),
	}
	db.BinLog.WriteEntry(bin)
	return db.Set(key, value, ttl)
}

func HandleGetCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: get [key]")
	}
	key := []byte(parts[1])
	value, err := db.Get(key)
	if err != nil {
		return err
	}
	fmt.Println("Value:", string(value))
	return nil
}

func HandleAppendCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 3 {
		return errors.New("Usage: append [key] [value]")
	}
	key := []byte(parts[1])
	value := []byte(parts[2])
	return db.Append(key, value)
}

func HandleStrlenCommand(parts []string, db *database.XcDB) error {
	if len(parts) != 2 {
		return errors.New("Usage: strlen [key]")
	}
	key := []byte(parts[1])
	valueLen, err := db.Strlen(key)
	if err != nil {
		return err
	}
	fmt.Println("Value length:", valueLen)
	return nil
}
