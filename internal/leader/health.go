package leader

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "distributed-disk-register-with-grpc/proto/family"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartHealthChecker(registry *NodeRegistry, self *pb.NodeInfo) {
	go func() {
		for {
			time.Sleep(10 * time.Second)

			for _, n := range registry.Snapshot() {
				if n.Host == self.Host && n.Port == self.Port {
					continue
				}

				addr := fmt.Sprintf("%s:%d", n.Host, n.Port)
				conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					log.Printf("Node dead: %s:%d\n", n.Host, n.Port)
					registry.Remove(n)
					continue
				}

				client := pb.NewFamilyServiceClient(conn)
				_, err = client.GetFamily(context.Background(), &pb.Empty{})
				if err != nil {
					log.Printf("Node unreachable: %s:%d\n", n.Host, n.Port)
					registry.Remove(n)
				}

				conn.Close()
			}
		}
	}()
}
