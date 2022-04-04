package service

import (
	"context"
	"fmt"
	pb "github.com/iman_task/go-service/genproto/collect"
	loggerPkg "github.com/iman_task/go-service/pkg/logger"
)

func (s *GoService) CollectPosts(ctx context.Context, req *pb.CollectPostsRequest) (*pb.CollectPostsResponse, error) {

	err := s.storage.Collect().CollectPostsStart()
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to update post collection status from db"), loggerPkg.Error(err))
		return nil, err
	}

	return &pb.CollectPostsResponse{
		Errors: nil,
		Code:   0,
	}, nil
}

func (s *GoService) CheckStatus(ctx context.Context, req *pb.CheckStatusRequest) (*pb.CheckStatusResponse, error) {

	status, err := s.storage.Collect().CheckFinished()
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to check post collection status from db"), loggerPkg.Error(err))
		return nil, err
	}

	return &pb.CheckStatusResponse{
		Status: status,
		Errors: nil,
		Code:   0,
	}, nil
}
