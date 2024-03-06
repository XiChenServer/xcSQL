package grpc

import (
	"SQL/api/db_hash"
	"SQL/api/db_set"
	"SQL/internal/database"
)

type Server struct {
	DB *database.XcDB // 假设你的数据库实例叫做 XcDB
	db_hash.UnimplementedHashDatabaseServer
	db_set.UnimplementedSetDatabaseServer
}
