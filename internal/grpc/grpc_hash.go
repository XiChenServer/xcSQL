package grpc

import (
	"SQL/api/db_hash"
	"SQL/internal/database"
	"context"
)

type Server struct {
	DB *database.XcDB // 假设你的数据库实例叫做 XcDB
	db_hash.UnimplementedHashDatabaseServer
}

func (s *Server) HSet(ctx context.Context, req *db_hash.HSetRequest) (*db_hash.HSetResponse, error) {
	err := s.DB.HSet(req.Key, req.Values, req.Ttl...)
	if err != nil {
		return &db_hash.HSetResponse{Success: false}, err
	}

	return &db_hash.HSetResponse{Success: true}, nil
}

func (s *Server) HGet(ctx context.Context, req *db_hash.HGetRequest) (*db_hash.HGetResponse, error) {
	value, err := s.DB.HGet(req.Key, req.Field)
	if err != nil {
		return nil, err
	}
	return &db_hash.HGetResponse{Value: value}, nil
}

func (s *Server) HGetAll(ctx context.Context, req *db_hash.HGetAllRequest) (*db_hash.HGetAllResponse, error) {
	values, err := s.DB.HGETALL(req.Key)
	if err != nil {
		return nil, err
	}
	return &db_hash.HGetAllResponse{Values: values}, nil
}

func (s *Server) HDel(ctx context.Context, req *db_hash.HDelRequest) (*db_hash.HDelResponse, error) {
	var success bool
	err := s.DB.HDel(req.Key, req.Fields...)
	if err != nil {
		return nil, err
	}
	success = true
	return &db_hash.HDelResponse{Success: success}, nil
}

func (s *Server) HExists(ctx context.Context, req *db_hash.HExistsRequest) (*db_hash.HExistsResponse, error) {
	exists, err := s.DB.HExists(req.Key, req.Field)
	if err != nil {
		return nil, err
	}
	return &db_hash.HExistsResponse{Exists: exists}, nil
}

func (s *Server) HKeys(ctx context.Context, req *db_hash.HKeysRequest) (*db_hash.HKeysResponse, error) {
	keys, err := s.DB.HKeys(req.Key)
	if err != nil {
		return nil, err
	}
	return &db_hash.HKeysResponse{Keys: keys}, nil
}

func (s *Server) HVals(ctx context.Context, req *db_hash.HValsRequest) (*db_hash.HValsResponse, error) {
	values, err := s.DB.HVals(req.Key)
	if err != nil {
		return nil, err
	}
	return &db_hash.HValsResponse{Values: values}, nil
}

func (s *Server) HLen(ctx context.Context, req *db_hash.HLenRequest) (*db_hash.HLenResponse, error) {
	length, err := s.DB.HLen(req.Key)
	if err != nil {
		return nil, err
	}
	return &db_hash.HLenResponse{Length: int32(length)}, nil
}
