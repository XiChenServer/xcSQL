package grpc

import (
	"SQL/api/db_set"
	"context"
)

// 实现 SetDatabase 服务接口的 SAdd 方法
func (s *Server) SAdd(ctx context.Context, req *db_set.SAddRequest) (*db_set.SAddResponse, error) {
	err := s.DB.SAdd(req.Key, req.Members, req.Ttl...)
	if err != nil {
		return &db_set.SAddResponse{
			Success: false,
		}, err
	}
	return &db_set.SAddResponse{Success: true}, nil
}

// 实现 SetDatabase 服务接口的 SRem 方法
func (s *Server) SRem(ctx context.Context, req *db_set.SRemRequest) (*db_set.SRemResponse, error) {
	err := s.DB.SRem(req.Key, req.Members)
	if err != nil {
		return &db_set.SRemResponse{
			Success: false,
		}, err
	}
	return &db_set.SRemResponse{Success: true}, nil
}

// 实现 SetDatabase 服务接口的 SMembers 方法
func (s *Server) SMembers(ctx context.Context, req *db_set.SMembersRequest) (*db_set.SMembersResponse, error) {
	members, err := s.DB.SMembers(req.Key)
	if err != nil {
		return nil, err
	}
	return &db_set.SMembersResponse{Members: members}, nil
}

// 实现 SetDatabase 服务接口的 SIsMember 方法
func (s *Server) SIsMember(ctx context.Context, req *db_set.SIsMemberRequest) (*db_set.SIsMemberResponse, error) {
	exists, err := s.DB.SIsMember(req.Key, req.Member)
	if err != nil {
		return nil, err
	}
	return &db_set.SIsMemberResponse{Exists: exists}, nil
}

// 实现 SetDatabase 服务接口的 SCard 方法
func (s *Server) SCard(ctx context.Context, req *db_set.SCardRequest) (*db_set.SCardResponse, error) {
	length, err := s.DB.SCard(req.Key)
	if err != nil {
		return nil, err
	}
	return &db_set.SCardResponse{Length: int32(length)}, nil
}
