package main

import (
	"SQL/api/db_hash"
	"SQL/internal/database"
	"SQL/internal/lsm"
	"SQL/internal/storage"
	"time"

	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	db_grpc "SQL/internal/grpc"
	"google.golang.org/grpc"
)

func main() {
	// 初始化日志记录器
	//logs.InitLogger()
	// 获取当前工作目录
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}
	log.Println("当前工作目录:", pwd)
	// 设置 TCP 监听器
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("无法监听: %v", err)
	}
	defer lis.Close()

	// 创建 gRPC 服务器
	s := grpc.NewServer()

	// 初始化哈希数据库
	db := database.NewConnectionPool(2, 10, 30*time.Minute)

	driver, err := db.GetConnection("123")
	if err != nil {
		log.Fatalf("无法获取数据库连接: %v", err)
	}
	defer db.ReleaseConnection(driver) // 确保释放连接
	// 将哈希数据库服务器注册到 gRPC 服务器
	db_hash.RegisterHashDatabaseServer(s, &db_grpc.Server{
		DB: driver.DB,
	})

	// 打印服务器启动消息
	log.Println("服务器启动在 :50051")

	// 在单独的 goroutine 中开始处理传入的 gRPC 请求
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("无法提供服务: %v", err)
		}
	}()

	// 设置信号通道以捕获终止信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待终止信号
	sig := <-sigChan
	log.Printf("收到信号 %v，开始优雅关闭...\n", sig)

	// 优雅关闭 gRPC 服务器
	s.GracefulStop()

	// 在退出时保存活动数据到磁盘并将磁盘数据打印到文件中以供 LSM 树使用
	saveAndPrintDiskData(driver.DB.Lsm)

	// 将存储管理器配置保存到文件
	storage.SaveStorageManager(driver.DB.StorageManager, "../../data/testdata/lsm_tree/config.txt")
	log.Println("服务器优雅关闭")
}

func saveAndPrintDiskData(lsmMap *map[uint16]*lsm.LSMTree) {
	for _, lsm := range *lsmMap {
		lsm.SaveActiveToDiskOnExit()
		lsm.PrintDiskDataToFile(string(lsm.LsmPath))
	}
}
