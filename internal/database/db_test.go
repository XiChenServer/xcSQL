package database

import (
	"fmt"
	"testing"
)

func Test_db(t *testing.T) {

	db := DBConnect("1")
	err := DBExit(db)
	if err != nil {
		fmt.Println(err)
	}

}

// 简单的测试数据可以存入
func TestDB_S(t *testing.T) {
	db := DBConnect("1")
	key := []byte(generateRandomKey())
	value := []byte(generateRandomKey())
	//fmt.Println(db.StorageManager.StoragePath)
	err := db.Set(key, value)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Insert ok")
	fmt.Println(string(key), string(value))

	err = DBExit(db)
	if err != nil {
		fmt.Println(err)
	}

}
