package main

import (
	"SQL/api/db_hash"
	"context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// 连接 gRPC 服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 创建一个 HashDatabase 客户端
	client := db_hash.NewHashDatabaseClient(conn)

	// 调用 HSet 方法
	hSetReq := &db_hash.HSetRequest{
		Key:    []byte("12323"),
		Values: map[string]string{"field1": "2323", "field2": "value2"},
		Ttl:    []uint64{3600}, // 设置 TTL 为 3600 秒
	}
	hSetResp, err := client.HSet(context.Background(), hSetReq)
	if err != nil {
		log.Fatalf("HSet failed: %v", err)
	}
	log.Printf("HSet Response: %v\n", hSetResp)

	// 调用 HGet 方法
	hGetReq := &db_hash.HGetRequest{
		Key:   []byte("12323"),
		Field: "field1",
	}
	hGetResp, err := client.HGet(context.Background(), hGetReq)
	if err != nil {
		log.Fatalf("HGet failed: %v", err)
	}
	log.Printf("HGet Response: %v\n", hGetResp)
}

//package main
//
//import (
//	"SQL/api/db_hash"
//	"google.golang.org/grpc"
//	"log"
//	"sync"
//	"time"
//)
//
//const (
//	concurrency = 100  // 并发数
//	requests    = 1000 // 每个 goroutine 的请求数
//)

//func main() {
//	// 连接 gRPC 服务器
//	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
//	if err != nil {
//		log.Fatalf("Failed to connect: %v", err)
//	}
//	defer conn.Close()
//
//	// 创建一个 HashDatabase 客户端
//	client := db_hash.NewHashDatabaseClient(conn)
//
//	// 使用 WaitGroup 来等待所有 goroutines 完成
//	var wg sync.WaitGroup
//
//	start := time.Now()
//
//	for i := 0; i < concurrency; i++ {
//		wg.Add(1)
//		go func(i int) {
//			defer wg.Done()
//			for j := 0; j < requests; j++ {
//				hSetReq := &db_hash.HSetRequest{
//					Key:    []byte("key" + string(i) + "_" + string(j)),
//					Values: map[string]string{"field1": "value1", "field2": "value2"},
//					Ttl:    []uint64{3600},
//				}
//				_, err := client.HSet(context.Background(), hSetReq)
//				if err != nil {
//					log.Printf("HSet failed for key %s: %v", hSetReq.Key, err)
//				}
//			}
//		}(i)
//	}
//
//	// 等待所有的 goroutines 完成
//	wg.Wait()
//
//	elapsed := time.Since(start)
//	log.Printf("Completed %d HSet requests with %d concurrency in %s", concurrency*requests, concurrency, elapsed)
//}

// package main
//
// import (
//
//	"SQL/api/db_hash"
//	"context"
//	"log"
//	"sync"
//	"time"
//
//	"google.golang.org/grpc"
//
// )
//
// const (
//
//	concurrency = 100  // 并发数
//	requests    = 1000 // 每个 goroutine 的请求数
//
// )
//
//	func main() {
//		// 连接 gRPC 服务器
//		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
//		if err != nil {
//			log.Fatalf("Failed to connect: %v", err)
//		}
//		defer conn.Close()
//
//		// 创建一个 HashDatabase 客户端
//		client := db_hash.NewHashDatabaseClient(conn)
//
//		// 使用 WaitGroup 来等待所有 goroutines 完成
//		var wg sync.WaitGroup
//
//		start := time.Now()
//
//		for i := 0; i < concurrency; i++ {
//			wg.Add(1)
//			go func(i int) {
//				defer wg.Done()
//				for j := 0; j < requests; j++ {
//					hGetReq := &db_hash.HGetRequest{
//						Key:   []byte("key" + string(i) + "_" + string(j)),
//						Field: "field1",
//					}
//					_, err := client.HGet(context.Background(), hGetReq)
//					if err != nil {
//						log.Printf("HGet failed for key %s: %v", hGetReq.Key, err)
//					}
//				}
//			}(i)
//		}
//
//		// 等待所有的 goroutines 完成
//		wg.Wait()
//
//		elapsed := time.Since(start)
//		log.Printf("Completed %d HGet requests with %d concurrency in %s", concurrency*requests, concurrency, elapsed)
//	}
//
// package main
const (
	concurrency = 1 // 并发数
	requests    = 1 // 每个 goroutine 的请求数
)

//func main() {
//	// 连接 gRPC 服务器
//	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
//	if err != nil {
//		log.Fatalf("Failed to connect: %v", err)
//	}
//	defer conn.Close()
//
//	// 创建一个 HashDatabase 客户端
//	client := db_hash.NewHashDatabaseClient(conn)
//
//	// 使用 WaitGroup 来等待所有 goroutines 完成
//	var wg sync.WaitGroup
//
//	start := time.Now()
//
//	// 启动写操作的 goroutines
//	for i := 0; i < concurrency; i++ {
//		wg.Add(1)
//		go func(i int) {
//			defer wg.Done()
//			for j := 0; j < requests; j++ {
//				hSetReq := &db_hash.HSetRequest{
//					Key:    []byte("1"),
//					Values: map[string]string{"field1": "value1", "field2": "value2"},
//					Ttl:    []uint64{3600},
//				}
//				_, err := client.HSet(context.Background(), hSetReq)
//				if err != nil {
//					log.Printf("HSet failed for key %s: %v", hSetReq.Key, err)
//				}
//			}
//		}(i)
//	}
//
//	// 启动读操作的 goroutines
//	for i := 0; i < concurrency; i++ {
//		wg.Add(1)
//		go func(i int) {
//			defer wg.Done()
//			for j := 0; j < requests; j++ {
//				hGetReq := &db_hash.HGetRequest{
//					Key:   []byte("1"),
//					Field: "field1",
//				}
//				_, err := client.HGet(context.Background(), hGetReq)
//				if err != nil {
//					log.Printf("HGet failed for key %s: %v", hGetReq.Key, err)
//				}
//			}
//		}(i)
//	}
//
//	// 等待所有的 goroutines 完成
//	wg.Wait()
//
//	elapsed := time.Since(start)
//	log.Printf("Completed %d HSet and %d HGet requests with %d concurrency in %s",
//		concurrency*requests, concurrency*requests, concurrency, elapsed)
//}
