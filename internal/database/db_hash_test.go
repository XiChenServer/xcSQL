package database

import (
	"math/rand"
	"reflect"
	"testing"
)

func Test_MapToBytesAndBytesToMap(t *testing.T) {
	// 定义一个 map[string]string
	value := generateRandomMap(4, 5)
	// 将 map[string]string 转换为 []byte
	bytes, err := mapToBytes(value)
	if err != nil {
		t.Errorf("mapToBytes returned an error: %v", err)
	}

	// 将 []byte 转换为 map[string]string
	m, err := bytesToMap(bytes)
	if err != nil {
		t.Errorf("bytesToMap returned an error: %v", err)
	}

	// 检查原始 map 和解析出来的 map 是否相等
	if !reflect.DeepEqual(value, m) {
		t.Errorf("mapToBytes and bytesToMap result mismatch: %v != %v", value, m)
	}
}

// 随机生成指定长度的字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// 生成随机的 map[string]string 数据
func generateRandomMap(length, size int) map[string]string {
	randomMap := make(map[string]string)
	for i := 0; i < size; i++ {
		key := randomString(length)
		value := randomString(length)
		randomMap[key] = value
	}
	return randomMap
}
