package service

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/goodguy-project/goodguy-crawl/v2/handler"
	"github.com/goodguy-project/goodguy-crawl/v2/proto"
)

type GrpcServer struct {
	proto.UnimplementedGoodguyCrawlServiceServer
}

func (GrpcServer) GetContestRecord(ctx context.Context, req *proto.GetContestRecordRequest) (*proto.GetContestRecordResponse, error) {
	return handler.GetContestRecord(ctx, req)
}

func (GrpcServer) GetSubmitRecord(ctx context.Context, req *proto.GetSubmitRecordRequest) (*proto.GetSubmitRecordResponse, error) {
	return handler.GetSubmitRecord(ctx, req)
}

func (GrpcServer) GetRecentContest(ctx context.Context, req *proto.GetRecentContestRequest) (*proto.GetRecentContestResponse, error) {
	return handler.GetRecentContest(ctx, req)
}

func (GrpcServer) GetDailyQuestion(ctx context.Context, req *proto.GetDailyQuestionRequest) (*proto.GetDailyQuestionResponse, error) {
	return handler.GetDailyQuestion(ctx, req)
}

func RunGrpcService() {
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:9851"))
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	proto.RegisterGoodguyCrawlServiceServer(server, GrpcServer{})
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
