package leader

import (
	"context"
	"distributed-disk-register-with-grpc/internal/logger"
	"distributed-disk-register-with-grpc/internal/node"
	pb "distributed-disk-register-with-grpc/proto/family"

	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendSetToMember(node *pb.NodeInfo, id int, text string) error {

	address := fmt.Sprintf("%s:%d", node.Host, node.Port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewFamilyServiceClient(conn)

	_, err = client.Store(context.Background(), &pb.StoredMessage{Id: int32(id), Text: text})
	if err != nil {
		logger.TCPEvent("Failed to send SET to follower %s:%d for ID %d: %v", node.Host, node.Port, id, err)
		return err
	}

	fmt.Printf("Sent SET to member %s:%d for ID %d\n", node.Host, node.Port, id)
	logger.TCPEvent("SET request confirmed on follower %s:%d for ID %d", node.Host, node.Port, id)
	return nil

}

func SendGetToMember(node *pb.NodeInfo, id int, registry *node.Registry) (string, error) {
	address := fmt.Sprintf("%s:%d", node.Host, node.Port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Connection failed, removing node from registry:", address)
		logger.TCPEvent("GET request failed to connect to follower %s:%d for ID %d: %v", node.Host, node.Port, id, err)
		registry.RemoveByAddress(node.Host, node.Port)
		return "", err

	}
	defer conn.Close()
	client := pb.NewFamilyServiceClient(conn)

	response, err := client.Retrieve(context.Background(), &pb.MessageId{Id: int32(id)})
	if err != nil {
		logger.TCPEvent("Failed to retrieve from follower %s:%d for ID %d: %v", node.Host, node.Port, id, err)
		return "", err
	}

	logger.TCPEvent("GET request confirmed on follower %s:%d for ID %d", node.Host, node.Port, id)
	return response.Text, nil

}
