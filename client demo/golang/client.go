package main

import (
	"context"
	"fmt"
	"time"

	"goodguy_crawl_client_demo/pb"

	"google.golang.org/grpc"
)

func GetClient() (*grpc.ClientConn, *pb.CrawlServiceClient, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", "localhost", 9851),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return conn, nil, err
	}
	client := pb.NewCrawlServiceClient(conn)
	return conn, &client, nil
}

func main() {
	conn, client, err := GetClient()
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("close grpc connection fail, err=%v\n", err)
		}
	}(conn)
	if err != nil {
		fmt.Printf("get grpc client fail, err=%v\n", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	platform := "codeforces"
	rsp, err := (*client).GetRecentContest(ctx, &pb.GetRecentContestRequest{
		Platform: platform,
	})
	if err != nil {
		fmt.Printf("get recent contest failed. platform=%s, err=%s\n", platform, err.Error())
		return
	}
	fmt.Println(rsp)
}
