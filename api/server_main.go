package main

import (
	"SQL/api/db_hash"
	"SQL/internal/database"
	"SQL/internal/model"
	"SQL/internal/storage"
	"SQL/logs"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	db *database.XcDB // 假设你的数据库实例叫做 XcDB
	db_hash.UnimplementedHashDatabaseServer
}

func (s *server) HSet(ctx context.Context, req *db_hash.HSetRequest) (*db_hash.HSetResponse, error) {
	logs.InitLogger()
	s.db = database.NewXcDB()
	//dataFilePath := "../../data/testdata/lsm_tree/test1.txt"
	lsmMap := *s.db.Lsm
	lsmType := lsmMap[model.XCDB_Hash]
	err := s.db.HSet(req.Key, req.Values, req.Ttl...)
	if err != nil {
		return &db_hash.HSetResponse{Success: false}, err
	}
	lsmType.SaveActiveToDiskOnExit()
	lsmType.PrintDiskDataToFile(string(lsmType.LsmPath))
	storage.SaveStorageManager(s.db.StorageManager, "../../data/testdata/lsm_tree/config.txt")
	return &db_hash.HSetResponse{Success: true}, nil
}

func (s *server) HGet(ctx context.Context, req *db_hash.HGetRequest) (*db_hash.HGetResponse, error) {
	value, err := s.db.HGet(req.Key, req.Field)
	if err != nil {
		return nil, err
	}
	return &db_hash.HGetResponse{Value: value}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	db_hash.RegisterHashDatabaseServer(s, &server{})

	log.Println("Server started at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
