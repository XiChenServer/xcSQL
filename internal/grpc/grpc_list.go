package grpc

import (
	"SQL/api/db_list"
	"context"
)

// RPUSH 实现了 gRPC 的 RPUSH 方法
func (s *Server) RPUSH(ctx context.Context, req *db_list.RPUSHRequest) (*db_list.RPUSHResponse, error) {
	err := s.DB.RPUSH(req.Key, req.Values, req.Ttl...)
	if err != nil {
		return &db_list.RPUSHResponse{
			Success: false,
		}, err
	}
	return &db_list.RPUSHResponse{Success: true}, nil
}

// LPUSH 实现了 gRPC 的 LPUSH 方法
func (s *Server) LPUSH(ctx context.Context, req *db_list.LPUSHRequest) (*db_list.LPUSHResponse, error) {
	err := s.DB.LPUSH(req.Key, req.Values, req.Ttl...)
	if err != nil {
		return &db_list.LPUSHResponse{
			Success: false,
		}, err
	}
	return &db_list.LPUSHResponse{Success: true}, nil
}

// LRANGE 实现了 gRPC 的 LRANGE 方法
func (s *Server) LRANGE(ctx context.Context, req *db_list.LRANGERequest) (*db_list.LRANGEResponse, error) {
	data, err := s.DB.LRANGE(req.Key, int(req.Left), int(req.Right))
	if err != nil {
		return nil, err
	}
	return &db_list.LRANGEResponse{Values: data}, nil
}

// LINDEX 实现了 gRPC 的 LINDEX 方法
func (s *Server) LINDEX(ctx context.Context, req *db_list.LINDEXRequest) (*db_list.LINDEXResponse, error) {
	data, err := s.DB.LINDEX(req.Key, int(req.Index))
	if err != nil {
		return nil, err
	}
	return &db_list.LINDEXResponse{Value: data}, nil
}

// LPOP 实现了 gRPC 的 LPOP 方法
func (s *Server) LPOP(ctx context.Context, req *db_list.LPOPRequest) (*db_list.LPOPResponse, error) {
	data, err := s.DB.LPOP(req.Key)
	if err != nil {
		return nil, err
	}
	return &db_list.LPOPResponse{Value: data}, nil
}

// RPOP 实现了 gRPC 的 RPOP 方法
func (s *Server) RPOP(ctx context.Context, req *db_list.RPOPRequest) (*db_list.RPOPResponse, error) {
	data, err := s.DB.RPOP(req.Key)
	if err != nil {
		return nil, err
	}
	return &db_list.RPOPResponse{Value: data}, nil
}

// LLEN 实现了 gRPC 的 LLEN 方法
func (s *Server) LLEN(ctx context.Context, req *db_list.LLENRequest) (*db_list.LLENResponse, error) {
	length, err := s.DB.LLEN(req.Key)
	if err != nil {
		return nil, err
	}
	return &db_list.LLENResponse{Length: int32(length)}, nil
}
