package bloomfilter

import (
	"strconv"
	"sync"
	"testing"
)

func TestLocalBloomService(t *testing.T) {
	// 初始化加密器
	encryptor := NewEncryptor()

	// 创建布隆过滤器实例
	m, k := int32(1000), int32(3)
	bloomService := NewLocalBloomService(m, k, encryptor)

	// 定义测试数据
	testCases := []struct {
		value    string
		expected bool
	}{
		{"apple", true},
		{"banana", true},
		{"cherry", false}, // 未加入布隆过滤器
	}

	// 添加前两个值到布隆过滤器
	bloomService.Set("apple")
	bloomService.Set("banana")

	// 逐个测试
	for _, tc := range testCases {
		if bloomService.Exist(tc.value) != tc.expected {
			t.Errorf("Exist(%s) = %v; expected %v", tc.value, bloomService.Exist(tc.value), tc.expected)
		}
	}
}

func TestLocalBloomService_Concurrency(t *testing.T) {
	// 初始化布隆过滤器
	m := int32(1024) // 位数组大小
	k := int32(4)    // 哈希函数的数量
	encryptor := NewEncryptor()
	bloom := NewLocalBloomService(m, k, encryptor)

	// 待测试的并发数量
	concurrency := 1000
	var wg sync.WaitGroup

	// 设置一些数据
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(val string) {
			defer wg.Done()
			bloom.Set(val)
		}(strconv.Itoa(i))
	}

	wg.Wait()

	// 测试是否存在
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(val string) {
			defer wg.Done()
			exists := bloom.Exist(val)
			if !exists {
				t.Errorf("Value %s should exist in Bloom filter but it does not", val)
			}
		}(strconv.Itoa(i))
	}

	wg.Wait()
}
