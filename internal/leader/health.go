package leader

import (
	"context"
	"log"
	"time"

	pb "distributed-disk-register-with-grpc/proto/family"

	"google.golang.org/grpc"
)

func StartHealthChecker(registry *NodeRegistry, self *pb.NodeInfo) {
	go func() {
		for {
			time.Sleep(10 * time.Second)

			for _, n := range registry.Snapshot() {
				if n.Host == self.Host && n.Port == self.Port {
					continue
				}

				addr := n.Host + ":" + string(rune(n.Port))
				conn, err := grpc.Dial(addr, grpc.WithInsecure())
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
