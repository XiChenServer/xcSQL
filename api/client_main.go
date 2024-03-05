package main

import (
	"SQL/api/db_hash"
	"context"
	"log"

	"google.golang.org/grpc"
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
		Key:    []byte("123"),
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
		Key:   []byte("123"),
		Field: "field1",
	}
	hGetResp, err := client.HGet(context.Background(), hGetReq)
	if err != nil {
		log.Fatalf("HGet failed: %v", err)
	}
	log.Printf("HGet Response: %v\n", hGetResp)
}
