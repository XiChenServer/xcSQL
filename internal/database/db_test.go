package database

import (
	"fmt"
	"math/rand"
	"testing"
)

func Test_db(t *testing.T) {

	db := DBConnect("1")
	fmt.Println("------------------------------------------")
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
	fmt.Println(string(db.StorageManager.StoragePath))
	fmt.Println("Insert ok")
	fmt.Println(string(key), string(value))

	err = DBExit(db)
	if err != nil {
		fmt.Println(err)
	}

}

// generateRandomKey 生成随机键值
func generateRandomKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLen := 10
	b := make([]byte, keyLen)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
