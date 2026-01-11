package node

import (
	"context"
	pb "distributed-disk-register-with-grpc/proto/family"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func JoinCluster(host string, leaderPort int32, nodeInfo *pb.NodeInfo) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, leaderPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed to connect to leader: %v\n", err)
		return
	}
	defer conn.Close()

	client := pb.NewFamilyServiceClient(conn)
	ctx := context.Background()

	client.Join(ctx, nodeInfo)
}
