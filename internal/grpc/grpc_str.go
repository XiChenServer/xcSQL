package grpc

import (
	"SQL/api/db_str"
	"context"
	"errors"
)

// Set 实现了 XcDB.Set RPC 方法
func (s *Server) Set(ctx context.Context, req *db_str.SetRequest) (*db_str.SetResponse, error) {
	err := s.DB.Set(req.Key, req.Value, req.Ttl...)
	if err != nil {
		return &db_str.SetResponse{ErrorMessage: err.Error()}, nil
	}
	return &db_str.SetResponse{}, nil
}

// Get 实现了 XcDB.Get RPC 方法
func (s *Server) Get(ctx context.Context, req *db_str.GetRequest) (*db_str.GetResponse, error) {
	value, err := s.DB.Get(req.Key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, errors.New("Not found")
	}
	return &db_str.GetResponse{Value: value}, nil
}

// Strlen 实现了 XcDB.Strlen RPC 方法
func (s *Server) Strlen(ctx context.Context, req *db_str.StrlenRequest) (*db_str.StrlenResponse, error) {
	length, err := s.DB.Strlen(req.Key)
	if err != nil {
		return nil, err
	}
	return &db_str.StrlenResponse{Length: length}, nil
}

// Append 实现了 XcDB.Append RPC 方法
func (s *Server) Append(ctx context.Context, req *db_str.AppendRequest) (*db_str.AppendResponse, error) {
	err := s.DB.Append(req.Key, req.Value)
	if err != nil {
		return &db_str.AppendResponse{ErrorMessage: err.Error()}, nil
	}
	return &db_str.AppendResponse{}, nil
}
