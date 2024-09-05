package alv_tree

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

// 基准测试插入操作的性能
func BenchmarkInsert(b *testing.B) {
	tree := NewAlvNode()
	b.ResetTimer()

	var wg sync.WaitGroup
	const numGoroutines = 10
	const numOperations = 10000
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := []byte(strconv.Itoa(i) + "-" + strconv.Itoa(j))
				node := nodeInfo{
					Key: key,
					TTL: time.Hour,
				}
				mu.Lock()
				err := tree.Insert(key, node)
				if err != nil {
					return
				}
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
}

// 基准测试删除操作的性能
func BenchmarkDelete(b *testing.B) {
	tree := NewAlvNode()
	// 插入测试数据
	for i := 0; i < 100000; i++ {
		key := []byte(strconv.Itoa(i))
		node := nodeInfo{
			Key: key,
			TTL: time.Hour,
		}
		err := tree.Insert(key, node)
		if err != nil {
			return
		}
	}
	b.ResetTimer()

	var wg sync.WaitGroup
	const numGoroutines = 10
	const numOperations = 10000
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := []byte(strconv.Itoa(i) + "-" + strconv.Itoa(j))
				mu.Lock()
				err := tree.Delete(key)
				if err != nil {
					return
				}
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
}

// 基准测试搜索操作的性能
func BenchmarkSearch(b *testing.B) {
	tree := NewAlvNode()
	// 插入测试数据
	for i := 0; i < 100000; i++ {
		key := []byte(strconv.Itoa(i))
		node := nodeInfo{
			Key: key,
			TTL: time.Hour,
		}
		err := tree.Insert(key, node)
		if err != nil {
			return
		}
	}
	b.ResetTimer()

	var wg sync.WaitGroup
	const numGoroutines = 10
	const numOperations = 10000
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := []byte(strconv.Itoa(j)) // 搜索的 key 是连续的
				mu.Lock()
				tree.Search(key)
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
}
